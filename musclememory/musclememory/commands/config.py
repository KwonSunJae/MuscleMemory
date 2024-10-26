import click
import os
import json
from ..utils.config import get_global_config

@click.group(name='config')
def config_group():
    """전역 설정 관리"""
    pass

@config_group.command()
def show():
    """저장소 URL 설정 보기"""
    repo_url = get_global_config()
    if repo_url:
        print(f"\n저장소 URL: {repo_url}")
    else:
        print("저장된 저장소 URL이 없습니다.")

@config_group.command()
@click.argument('repo_url')
def update(repo_url):
    """저장소 URL 업데이트"""
    home = os.path.expanduser("~")
    config_path = os.path.join(home, ".muscle", "config.json")
    
    config = {"repository_url": repo_url}
    with open(config_path, 'w') as f:
        json.dump(config, f, indent=2)
    print(f"저장소 URL이 업데이트되었습니다: {repo_url}")
