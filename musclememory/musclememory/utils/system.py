import os
import uuid
from .crypto import encrypt_mac_address, decrypt_mac_address

def get_mac_address(encrypt=True):
    """시스템의 첫 번째 NIC의 MAC 주소를 반환"""
    mac = ':'.join(['{:02x}'.format((uuid.getnode() >> elements) & 0xff) 
                    for elements in range(0,8*6,8)][::-1])
    return encrypt_mac_address(mac) if encrypt else mac

def verify_mac_address(encrypted_mac):
    """현재 시스템의 MAC 주소와 암호화된 MAC 주소 비교"""
    try:
        decrypted_mac = decrypt_mac_address(encrypted_mac)
        current_mac = get_mac_address(encrypt=False)
        return decrypted_mac == current_mac
    except:
        return False

def get_workspace_dir():
    """작업 디렉토리 경로 반환"""
    home = os.path.expanduser("~")
    workspace = os.path.join(home, ".muscle", "workspace")
    os.makedirs(workspace, exist_ok=True)
    return workspace
