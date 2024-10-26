import pytest
import os
import json
from musclememory.utils.system import get_workspace_dir, get_mac_address
from musclememory.utils.config import save_global_config, get_global_config
from musclememory.utils.crypto import encrypt_mac_address, decrypt_mac_address

def test_workspace_dir(mock_home_dir):
    """workspace 디렉토리 생성 테스트"""
    workspace = get_workspace_dir()
    assert os.path.exists(workspace)
    assert workspace.endswith(".muscle/workspace")

def test_mac_address():
    """MAC 주소 형식 테스트"""
    mac = get_mac_address(encrypt=False)
    assert len(mac) == 17  # XX:XX:XX:XX:XX:XX 형식
    assert mac.count(":") == 5

def test_crypto():
    """MAC 주소 암호화/복호화 테스트"""
    original = "00:11:22:33:44:55"
    encrypted = encrypt_mac_address(original)
    decrypted = decrypt_mac_address(encrypted)
    assert original == decrypted
    assert encrypted != original

def test_config(mock_home_dir):
    """전역 설정 저장/로드 테스트"""
    test_url = "https://github.com/test/repo"
    test_project = "test_project"
    
    # 설정이 없을 때
    assert get_global_config() is None
    
    # 설정 저장 후
    save_global_config(test_url, test_project)
    assert get_global_config() == test_url
