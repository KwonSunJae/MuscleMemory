import click
import os
import subprocess
import time

@click.group()
def cli():
    """MuscleMemory CLI"""
    pass

@cli.command()
@click.argument('repo_url')
def enroll(repo_url):
    """새 프로젝트나 사용자를 Muscle에 등록"""
    print(f"새 프로젝트를 {repo_url}에 등록 중...")

    # Git 저장소 클론 또는 최신 상태로 유지
    try:
        if not os.path.exists('repo'):
            subprocess.run(["git", "clone", repo_url, "repo"], check=True)
            print(f"{repo_url} 경로에 Git 저장소가 클론되었습니다.")
        else:
            subprocess.run(["git", "pull"], cwd="repo", check=True)
            print("Git 저장소가 최신 상태로 업데이트되었습니다.")
    except subprocess.CalledProcessError as e:
        print(f"Git 작업 중 오류 발생: {e}")
    
    print("등록 완료.")

LOCK_FILE_PATH = "repo/tfstate.lock"

@cli.command()
def lock():
    """잠금 파일 생성"""
    if os.path.exists(LOCK_FILE_PATH):
        print("잠금 파일이 이미 존재합니다. 다른 사용자가 작업 중입니다.")
        return False
    else:
        with open(LOCK_FILE_PATH, 'w') as lock_file:
            lock_file.write("잠금 파일 생성 시간: " + time.ctime())
        print("잠금 파일이 생성되었습니다.")
        return True

TFSTATE_FILE_PATH = "repo/tfstate.tfstate"

@cli.command()
def prepare_tfstate():
    """tfstate 파일을 준비하고 초기화"""
    if os.path.exists(TFSTATE_FILE_PATH):
        print("tfstate 파일이 이미 존재합니다. 최신 상태로 업데이트 중...")
        try:
            subprocess.run(["git", "pull"], cwd="repo", check=True)
            print("tfstate 파일이 최신 상태로 업데이트되었습니다.")
        except subprocess.CalledProcessError as e:
            print(f"Git 업데이트 중 오류 발생: {e}")
    else:
        print("tfstate 파일이 존재하지 않습니다. 기본 템플릿을 생성합니다.")
        os.makedirs("repo", exist_ok=True)  # 디렉토리가 없는 경우 생성
        with open(TFSTATE_FILE_PATH, 'w') as tfstate_file:
            tfstate_file.write('{}')  # 빈 JSON 객체로 초기화
        try:
            subprocess.run(["git", "add", "tfstate.tfstate"], cwd="repo", check=True)
            subprocess.run(["git", "commit", "-m", "초기 tfstate 파일 생성"], cwd="repo", check=True)
            subprocess.run(["git", "push"], cwd="repo", check=True)
            print("새로운 tfstate 파일이 Git에 커밋되었습니다.")
        except subprocess.CalledProcessError as e:
            print(f"Git 커밋 중 오류 발생: {e}")

if __name__ == "__main__":
    cli()
