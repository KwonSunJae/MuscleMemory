import click
import os
import git
import json
from datetime import datetime, timedelta
from ..utils.system import get_mac_address, get_workspace_dir

@click.command()
@click.argument('project_name')
def ready(project_name):
    """프로젝트의 .lock 파일을 통해 동시성 제어"""
    try:
        workspace = get_workspace_dir()
        project_path = os.path.join(workspace, project_name)
        
        if not os.path.exists(project_path):
            print(f"Error: {project_name} 프로젝트가 없습니다.")
            print(f"먼저 'muscle enroll {project_name}'을 실행해주세요.")
            return False

        current_owner = get_mac_address()
        expired_time = (datetime.now() + timedelta(minutes=30)).isoformat()
        
        inner_project_path = os.path.join(project_path, project_name)
        lock_file_path = os.path.join(inner_project_path, ".lock")
        
        repo = git.Repo(project_path)
        repo.git.checkout(project_name)
        repo.git.pull('origin', project_name)
            
        os.makedirs(inner_project_path, exist_ok=True)
        
        if os.path.exists(lock_file_path):
            with open(lock_file_path, 'r') as f:
                lock_data = json.load(f)
                
            stored_owner = lock_data.get('owner')
            stored_expired_time = datetime.fromisoformat(lock_data.get('expired_time'))
            
            current_owner = get_mac_address()  # 이미 암호화된 상태
            
            if stored_owner == current_owner:
                if datetime.now() < stored_expired_time:
                    lock_data['expired_time'] = expired_time
                    with open(lock_file_path, 'w') as f:
                        json.dump(lock_data, f)
                    print(f"작업 시간이 {expired_time}까지 연장되었습니다.")
            else:
                if datetime.now() > stored_expired_time:
                    lock_data = {
                        'owner': current_owner,
                        'expired_time': expired_time
                    }
                    with open(lock_file_path, 'w') as f:
                        json.dump(lock_data, f)
                    print(f"만료된 잠금을 획득했습니다. {expired_time}까지 작업 가능합니다.")
                else:
                    raise Exception(f"현재 다른 사용자가 작업 중입니다. ({stored_expired_time}까지)")
        else:
            lock_data = {
                'owner': current_owner,
                'expired_time': expired_time
            }
            with open(lock_file_path, 'w') as f:
                json.dump(lock_data, f)
            print(f"잠금을 획득했습니다. {expired_time}까지 작업 가능합니다.")
        
        repo.index.add([os.path.relpath(lock_file_path, project_path)])
        repo.index.commit(f"Update lock file for {project_name}")
        repo.git.push('origin', project_name)
        
        return True
        
    except Exception as e:
        print(f"에러 발생: {str(e)}")
        return False
