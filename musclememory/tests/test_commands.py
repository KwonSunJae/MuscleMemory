from click.testing import CliRunner
from musclememory.cli import cli
import os
import git
import json
import pytest

def test_init_command(mock_home_dir, mock_repo):
    """init 명령어 테스트"""
    repo, repo_path = mock_repo
    runner = CliRunner()
    result = runner.invoke(cli, ['init', repo_path, 'test_project', 'terraform'])
    
    assert result.exit_code == 0
    assert "전역 설정에 저장소 URL이 저장되었습니다" in result.output
    
    # 설정 파일 확인
    config_path = os.path.join(mock_home_dir, "config.json")
    with open(config_path, 'r') as f:
        config = json.load(f)
    assert config["repository_url"] == repo_path
    
    # 브랜치 확인
    workspace = os.path.join(mock_home_dir, "workspace")
    project_path = os.path.join(workspace, "test_project")
    repo = git.Repo(project_path)
    
    assert "main" in repo.heads
    assert "blank" in repo.heads
    assert "test_project" in repo.heads
    
    # main 브랜치 확인
    repo.heads.main.checkout()
    assert os.path.exists(os.path.join(project_path, "muscle.init"))
    
    # blank 브랜치 확인
    repo.heads.blank.checkout()
    # blank 브랜치의 최신 커밋 메시지 확인
    assert repo.head.commit.message.strip() == "Initialize blank branch"
    
    # project 브랜치 확인
    repo.heads.test_project.checkout()
    project_dir = os.path.join(project_path, "test_project")
    assert os.path.exists(project_dir)
    assert os.path.exists(os.path.join(project_dir, "config"))
    
    # config 파일 내용 확인
    with open(os.path.join(project_dir, "config"), 'r') as f:
        config_content = f.read()
        assert "project_type=terraform" in config_content
        assert "project_name=test_project" in config_content
        assert f"project_repo={repo_path}" in config_content

def test_ready_command(mock_home_dir, mock_repo):
    """ready 명령어 테스트"""
    repo, repo_path = mock_repo
    
    # 먼저 init 실행
    runner = CliRunner()
    runner.invoke(cli, ['init', repo_path, 'test_project', 'terraform'])
    
    # ready 명령 테스트
    result = runner.invoke(cli, ['ready', 'test_project'])
    assert result.exit_code == 0
    
    # lock 파일 확인
    workspace = os.path.join(mock_home_dir, "workspace")
    project_path = os.path.join(workspace, "test_project")
    lock_path = os.path.join(project_path, "test_project", ".lock")
    assert os.path.exists(lock_path)
    
    with open(lock_path, 'r') as f:
        lock_data = json.load(f)
        assert "owner" in lock_data
        assert "expired_time" in lock_data

def test_complete_command(mock_home_dir, mock_repo):
    """complete 명령어 테스트"""
    repo, repo_path = mock_repo
    runner = CliRunner()
    
    # init과 ready 먼저 실행
    runner.invoke(cli, ['init', repo_path, 'test_project', 'terraform'])
    runner.invoke(cli, ['ready', 'test_project'])
    
    # complete 명령 테스트
    result = runner.invoke(cli, ['complete', 'test_project'])
    assert result.exit_code == 0
    
    # lock 파일이 제거되었는지 확인
    workspace = os.path.join(mock_home_dir, "workspace")
    project_path = os.path.join(workspace, "test_project")
    lock_path = os.path.join(project_path, "test_project", ".lock")
    assert not os.path.exists(lock_path)

# 기존 테스트 코드 아래에 추가

def test_add_command(mock_home_dir, mock_repo):
    """add 명령어 테스트"""
    repo, repo_path = mock_repo
    runner = CliRunner()
    
    # 초기 설정
    runner.invoke(cli, ['init', repo_path, 'test_project', 'terraform'])
    runner.invoke(cli, ['ready', 'test_project'])
    
    # 테스트 파일 생성
    workspace = os.path.join(mock_home_dir, "workspace")
    project_path = os.path.join(workspace, "test_project")
    test_file_path = os.path.join(project_path, "test_project", "main.tf")
    
    with open(test_file_path, 'w') as f:
        f.write('resource "aws_instance" "example" {}')
    
    # add 명령 테스트
    result = runner.invoke(cli, ['add', 'test_project', 'main.tf'])
    assert result.exit_code == 0
    
    # Git 상태 확인
    repo = git.Repo(project_path)
    assert not repo.is_dirty()  # 변경사항이 모두 스테이징되었는지 확인
    assert "main.tf" in [item.a_path for item in repo.head.commit.tree.traverse()]

def test_add_with_pattern(mock_home_dir, mock_repo):
    """add 명령어 패턴 테스트"""
    repo, repo_path = mock_repo
    runner = CliRunner()
    
    # 초기 설정
    runner.invoke(cli, ['init', repo_path, 'test_project', 'terraform'])
    runner.invoke(cli, ['ready', 'test_project'])
    
    # 여러 테스트 파일 생성
    workspace = os.path.join(mock_home_dir, "workspace")
    project_path = os.path.join(workspace, "test_project")
    project_dir = os.path.join(project_path, "test_project")
    
    # .tf 파일들 생성
    files = ["main.tf", "variables.tf", "outputs.tf", "test.txt"]
    for file in files:
        with open(os.path.join(project_dir, file), 'w') as f:
            f.write(f'# Test content for {file}')
    
    # *.tf 패턴으로 add 테스트
    result = runner.invoke(cli, ['add', 'test_project', '*.tf'])
    assert result.exit_code == 0
    
    # Git 상태 확인
    repo = git.Repo(project_path)
    committed_files = [item.a_path for item in repo.head.commit.tree.traverse()]
    
    # .tf 파일들은 커밋되고 .txt 파일은 커밋되지 않았는지 확인
    assert all(f"test_project/{f}" in committed_files for f in files if f.endswith('.tf'))
    assert f"test_project/test.txt" not in committed_files

def test_enroll_command(mock_home_dir, mock_repo):
    """enroll 명령어 테스트"""
    repo, repo_path = mock_repo
    runner = CliRunner()
    
    # 초기 설정
    runner.invoke(cli, ['init', repo_path, 'test_project', 'terraform'])
    runner.invoke(cli, ['ready', 'test_project'])
    
    # 테스트 파일 생성 및 add
    workspace = os.path.join(mock_home_dir, "workspace")
    project_path = os.path.join(workspace, "test_project")
    test_file_path = os.path.join(project_path, "test_project", "main.tf")
    
    with open(test_file_path, 'w') as f:
        f.write('resource "aws_instance" "example" {}')
    
    runner.invoke(cli, ['add', 'test_project', 'main.tf'])
    
    # enroll 명령 테스트
    result = runner.invoke(cli, ['enroll', 'test_project'])
    assert result.exit_code == 0
    
    # Git 상태 확인
    repo = git.Repo(project_path)
    
    # 현재 브랜치가 test_project인지 확인
    assert repo.active_branch.name == "test_project"
    
    # 커밋 메시지 확인
    latest_commit = repo.head.commit
    assert "Enroll changes" in latest_commit.message
    
    # 파일이 커밋에 포함되어 있는지 확인
    assert "test_project/main.tf" in [item.a_path for item in latest_commit.tree.traverse()]

def test_enroll_with_empty_changes(mock_home_dir, mock_repo):
    """변경사항이 없을 때 enroll 명령어 테스트"""
    repo, repo_path = mock_repo
    runner = CliRunner()
    
    # 초기 설정
    runner.invoke(cli, ['init', repo_path, 'test_project', 'terraform'])
    runner.invoke(cli, ['ready', 'test_project'])
    
    # 변경사항 없이 enroll 실행
    result = runner.invoke(cli, ['enroll', 'test_project'])
    assert result.exit_code == 0
    assert "변경사항이 없습니다" in result.output
