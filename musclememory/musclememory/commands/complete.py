import click
import os
import git
from ..utils.system import get_workspace_dir

@click.command()
@click.argument('project_name')
def complete(project_name):
    """작업 완료 및 .lock 파일 제거"""
    try:
        workspace = get_workspace_dir()
        project_path = os.path.join(workspace, project_name)
        
        if not os.path.exists(project_path):
            print(f"Error: {project_name} 프로젝트가 없습니다.")
            return False
            
        inner_project_path = os.path.join(project_path, project_name)
        if not os.path.exists(inner_project_path):
            print(f"Error: {project_name} 프로젝트의 작업 디렉토리가 없습니다.")
            return False
            
        repo = git.Repo(project_path)
        repo.git.checkout(project_name)
        
        try:
            repo.git.pull('origin', project_name)
        except git.exc.GitCommandError as e:
            print(f"원격 저장소 동기화 중 오류: {e}")
            return False
            
        lock_file_path = os.path.join(inner_project_path, ".lock")
        
        try:
            repo.git.rm('-f', f"{project_name}/.lock")
            repo.index.commit(f"Remove lock file for {project_name}")
            repo.git.push('origin', project_name)
            print(f"{project_name}의 잠금이 해제되었습니다.")
        except git.exc.GitCommandError:
            if os.path.exists(lock_file_path):
                os.remove(lock_file_path)
                repo.index.add([os.path.relpath(lock_file_path, project_path)])
                repo.index.commit(f"Remove lock file for {project_name}")
                repo.git.push('origin', project_name)
                print(f"{project_name}의 잠금이 해제되었습니다.")
            else:
                print("이미 잠금이 해제되었습니다.")
        
        return True
        
    except Exception as e:
        print(f"에러 발생: {str(e)}")
        return False
