import click
import os
import git
import shutil
from ..utils.system import get_workspace_dir

@click.command()
@click.argument('project_name')
@click.argument('source_path')
def add(project_name, source_path):
    """프로젝트에 파일/디렉토리 추가 (로컬에만 저장)"""
    try:
        workspace = get_workspace_dir()
        project_path = os.path.join(workspace, project_name)
        inner_project_path = os.path.join(project_path, project_name)
        
        if not os.path.exists(project_path):
            print(f"Error: {project_name} 프로젝트가 없습니다.")
            return False
            
        if not os.path.exists(source_path):
            print(f"Error: {source_path}가 존재하지 않습니다.")
            return False
            
        dest_path = os.path.join(inner_project_path, os.path.basename(source_path))
        if os.path.isdir(source_path):
            shutil.copytree(source_path, dest_path)
        else:
            shutil.copy2(source_path, dest_path)
            
        repo = git.Repo(project_path)
        repo.git.checkout(project_name)
        repo.index.add([os.path.relpath(dest_path, project_path)])
        repo.index.commit(f"Add {os.path.basename(source_path)} to {project_name}")
        
        print(f"{source_path}가 {project_name}에 추가되었습니다. (로컬)")
        print(f"변경사항을 적용하려면 'muscle enroll {project_name}'을 실행하세요.")
        return True
        
    except Exception as e:
        print(f"에러 발생: {str(e)}")
        return False
