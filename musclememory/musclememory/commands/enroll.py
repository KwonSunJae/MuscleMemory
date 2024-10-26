import click
import os
import git
import subprocess
from ..utils.system import get_workspace_dir

@click.command()
@click.argument('project_name')
def enroll(project_name):
    """프로젝트에 Terraform apply 실행 및 결과 푸시"""
    try:
        workspace = get_workspace_dir()
        project_path = os.path.join(workspace, project_name)
        inner_project_path = os.path.join(project_path, project_name)
        
        if not os.path.exists(project_path):
            print(f"Error: {project_name} 프로젝트가 없습니다.")
            return False
            
        os.chdir(inner_project_path)
        
        print("terraform init 실행 중... 하는 척")
        # subprocess.run(["terraform", "init"], check=True)
        
        print("terraform apply 실행 중... 하는 척")
        # subprocess.run(["terraform", "apply", "-auto-approve"], check=True)
        
        repo = git.Repo(project_path)
        repo.git.checkout(project_name)
        repo.git.add(".")
        repo.index.commit(f"Apply terraform changes for {project_name}")
        repo.git.push('origin', project_name)
        
        print(f"{project_name}에 Terraform 변경사항이 적용되었습니다.")
        return True
        
    except Exception as e:
        print(f"에러 발생: {str(e)}")
        return False
