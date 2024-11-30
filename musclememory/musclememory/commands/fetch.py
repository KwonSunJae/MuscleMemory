import click
import os
import git
from ..utils.system import get_workspace_dir
from ..models.project import Project

@click.group()
def fetch():
    """IaC 파일 관리 명령어 그룹"""
    pass

@fetch.command()
@click.argument('project_name')
def latest(project_name):
    """현재 프로젝트의 최신 IaC 파일을 불러옴"""
    try:
        workspace = get_workspace_dir()
        project_path = os.path.join(workspace, project_name)
        
        # 프로젝트 존재 여부 확인
        if not os.path.exists(project_path):
            print(f"프로젝트 {project_name}을 찾을 수 없습니다.")
            return False
            
        repo = git.Repo(project_path)
        
        # 프로젝트 브랜치로 전환
        repo.git.checkout(project_name)
        
        # 원격 저장소에서 최신 변경사항 가져오기
        repo.git.fetch('origin', project_name)
        
        # 변경사항 있는지 확인
        if repo.is_dirty():
            print("저장되지 않은 로컬 변경사항이 있습니다. 변경사항을 처리한 후 다시 시도해주세요.")
            return False
            
        repo.git.reset('--hard', f'origin/{project_name}')
        print(f"{project_name}의 최신 IaC 파일을 가져왔습니다.")
        return True
        
    except git.exc.GitCommandError as e:
        print(f"Git 작업 중 오류 발생: {e}")
        return False
    except Exception as e:
        print(f"IaC 파일 불러오기 실패: {e}")
        return False

@fetch.command()
@click.argument('project_name')
@click.argument('version')
def version(project_name, version):
    """특정 버전의 IaC 파일을 불러옴"""
    try:
        workspace = get_workspace_dir()
        project_path = os.path.join(workspace, project_name)
        
        # 프로젝트 존재 여부 확인
        if not os.path.exists(project_path):
            print(f"프로젝트 {project_name}을 찾을 수 없습니다.")
            return False
            
        repo = git.Repo(project_path)
        
        # 프로젝트 브랜치로 전환
        repo.git.checkout(project_name)
        
        # 변경사항 있는지 확인
        if repo.is_dirty():
            print("저장되지 않은 로컬 변경사항이 있습니다. 변경사항을 처리한 후 다시 시도해주세요.")
            return False
            
        # 특정 버전으로 이동
        try:
            repo.git.checkout(version)
            print(f"{project_name}의 {version} 버전 IaC 파일을 가져왔습니다.")
            return True
        except git.exc.GitCommandError:
            print(f"버전 {version}을 찾을 수 없습니다.")
            return False
            
    except git.exc.GitCommandError as e:
        print(f"Git 작업 중 오류 발생: {e}")
        return False
    except Exception as e:
        print(f"특정 버전 IaC 파일 불러오기 실패: {e}")
        return False
