from cryptography.fernet import Fernet
import base64
import os

def get_or_create_key():
    """암호화 키 가져오기 또는 생성"""
    home = os.path.expanduser("~")
    key_path = os.path.join(home, ".muscle", "crypto.key")
    
    if os.path.exists(key_path):
        with open(key_path, 'rb') as f:
            return f.read()
    else:
        key = Fernet.generate_key()
        os.makedirs(os.path.dirname(key_path), exist_ok=True)
        with open(key_path, 'wb') as f:
            f.write(key)
        return key

def encrypt_mac_address(mac_address: str) -> str:
    """MAC 주소 암호화"""
    key = get_or_create_key()
    f = Fernet(key)
    encrypted = f.encrypt(mac_address.encode())
    return base64.urlsafe_b64encode(encrypted).decode()

def decrypt_mac_address(encrypted_mac: str) -> str:
    """MAC 주소 복호화"""
    key = get_or_create_key()
    f = Fernet(key)
    encrypted = base64.urlsafe_b64decode(encrypted_mac.encode())
    return f.decrypt(encrypted).decode()
