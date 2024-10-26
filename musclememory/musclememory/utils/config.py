import os
import json

def save_global_config(repo_url, project_name):
    """전역 설정 저장"""
    home = os.path.expanduser("~")
    muscle_dir = os.path.join(home, ".muscle")
    os.makedirs(muscle_dir, exist_ok=True)
    
    config_path = os.path.join(muscle_dir, "config.json")
    config = {}
    
    if os.path.exists(config_path):
        with open(config_path, 'r') as f:
            config = json.load(f)
    
    config["repository_url"] = repo_url
    
    with open(config_path, 'w') as f:
        json.dump(config, f, indent=2)

def get_global_config(project_name=None):
    """전역 설정에서 저장소 URL 가져오기"""
    home = os.path.expanduser("~")
    config_path = os.path.join(home, ".muscle", "config.json")
    
    if os.path.exists(config_path):
        with open(config_path, 'r') as f:
            config = json.load(f)
            return config.get("repository_url")
    return None

def get_repo_url(project_name):
    """프로젝트의 config 파일에서 repo_url을 가져옴"""
    config_path = os.path.join(project_name, project_name, "config")
    if os.path.exists(config_path):
        with open(config_path, 'r') as f:
            for line in f:
                if line.startswith("project_repo="):
                    return line.strip().split("=")[1]
    return None
