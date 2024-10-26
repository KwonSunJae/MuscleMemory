import click
import os
import git
from git.exc import GitCommandError, InvalidGitRepositoryError, NoSuchPathError
from ..utils.system import get_mac_address, get_workspace_dir
from ..utils.config import save_global_config

@click.command()
@click.argument('repo_url')
@click.argument('project_name')
@click.argument('config_type')
def init(repo_url, project_name, config_type):
    """muscle init으로 전체 초기화 작업 수행"""
    try:
        # 전역 설정 저장
        save_global_config(repo_url, project_name)
        print(f"전역 설정에 저장소 URL이 저장되었습니다.")

        # workspace 디렉토리 내에 프로젝트 생성
        workspace = get_workspace_dir()
        project_path = os.path.join(workspace, project_name)

        # 기존 디렉토리가 있다면 제거
        if os.path.exists(project_path):
            import shutil
            shutil.rmtree(project_path)

        # 빈 저장소 초기화
        os.makedirs(project_path)
        repo = git.Repo.init(project_path)
        repo.create_remote('origin', repo_url)
        
        try:
            # 원격 저장소에서 main 브랜치 가져오기 시도
            repo.git.fetch('origin', 'main')
            repo.git.checkout('main')
            print("기존 main 브랜치를 가져왔습니다.")
        except GitCommandError:
            # main 브랜치가 없으면 새로 생성
            print("main 브랜치가 없어 새로 생성합니다.")
            # muscle.init 파일 생성
            muscle_init_path = os.path.join(project_path, "muscle.init")
            mac_address = get_mac_address()
            timeout = 300
            with open(muscle_init_path, "w") as init_file:
                init_file.write(f"timeout={timeout}\n")
                init_file.write(f"mac_address={mac_address}\n")
                
            # main 브랜치 생성 및 커밋
            repo.index.add(["muscle.init"])
            repo.index.commit("Add muscle.init file")
            repo.git.branch('-M', 'main')
            repo.git.push('-u', 'origin', 'main')
            print("main 브랜치 생성 및 푸시 완료")
        
        # blank 브랜치 처리
        try:
            # 원격 blank 브랜치 가져오기 시도
            repo.git.fetch('origin', 'blank')
            repo.git.checkout('blank')
            print("기존 blank 브랜치를 가져왔습니다.")
        except GitCommandError:
            # blank 브랜치가 없으면 새로 생성
            print("blank 브랜치가 없어 새로 생성합니다.")
            repo.git.checkout('-b', 'blank')
            repo.index.commit("Initialize blank branch")
            repo.git.push('-u', 'origin', 'blank')
            print("blank 브랜치 생성 및 푸시 완료")
        
        # project 브랜치 처리
        try:
            # 원격 project 브랜치 가져오기 시도
            repo.git.fetch('origin', project_name)
            repo.git.checkout(project_name)
            print(f"기존 {project_name} 브랜치를 가져왔습니다.")
        except GitCommandError:
            # project 브랜치가 없으면 새로 생성
            print(f"{project_name} 브랜치가 없어 새로 생성합니다.")
            repo.git.checkout('-b', project_name)
            
            # 프로젝트 디렉토리 및 설정 파일 생성
            inner_project_path = os.path.join(project_path, project_name)
            os.makedirs(inner_project_path, exist_ok=True)
            config_path = os.path.join(inner_project_path, "config")
            with open(config_path, "w") as config_file:
                config_file.write(f"project_type={config_type}\n")
                config_file.write(f"project_name={project_name}\n")
                config_file.write(f"project_repo={repo_url}\n")
                
            repo.index.add([os.path.relpath(inner_project_path, project_path)])
            repo.index.commit(f"Initialize {project_name} project directory")
            repo.git.push('-u', 'origin', project_name)
            print(f"{project_name} 브랜치 생성 및 푸시 완료")
        
        print(f"저장소가 성공적으로 초기화되었습니다.")
        return True

    except Exception as e:
        print(f"Git 작업 중 오류 발생: {e}")
        return False
