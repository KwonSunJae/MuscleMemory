import pytest
import os
import git
from click.testing import CliRunner
from musclememory.cli import cli

def test_rollback_basic(mock_home_dir, mock_repo):
    """기본 롤백 테스트"""
    repo, repo_path = mock_repo
    runner = CliRunner()
    
    # 초기 설정
    runner.invoke(cli, ['init', repo_path, 'test_project', 'terraform'])
    
    workspace = os.path.join(mock_home_dir, "workspace")
    project_path = os.path.join(workspace, "test_project")
    test_file_path = os.path.join(project_path, "test_project", "main.tf")
    
    repo = git.Repo(project_path)
    
    # 첫 번째 버전
    with open(test_file_path, 'w') as f:
        f.write('resource "aws_instance" "example" { instance_type = "t2.micro" }')
    repo.git.add('test_project/main.tf')
    repo.git.commit('-m', 'Version 1')
    first_commit = repo.head.commit.hexsha
    
    # 두 번째 버전
    with open(test_file_path, 'w') as f:
        f.write('resource "aws_instance" "example" { instance_type = "t2.small" }')
    repo.git.add('test_project/main.tf')
    repo.git.commit('-m', 'Version 2')
    
    # 롤백 테스트 (자동 No 응답)
    result = runner.invoke(cli, ['rollback', 'test_project', first_commit], input='n\n')
    assert result.exit_code == 0
    assert f"버전 {first_commit}으로 롤백했습니다" in result.output
    
    # 파일 내용 확인
    with open(test_file_path, 'r') as f:
        content = f.read()
    assert 't2.micro' in content

def test_rollback_with_backup(mock_home_dir, mock_repo):
    """백업 옵션을 사용한 롤백 테스트"""
    repo, repo_path = mock_repo
    runner = CliRunner()
    
    # 초기 설정
    runner.invoke(cli, ['init', repo_path, 'test_project', 'terraform'])
    
    workspace = os.path.join(mock_home_dir, "workspace")
    project_path = os.path.join(workspace, "test_project")
    test_file_path = os.path.join(project_path, "test_project", "main.tf")
    
    repo = git.Repo(project_path)
    
    # 첫 번째 버전
    with open(test_file_path, 'w') as f:
        f.write('Version 1 content')
    repo.git.add('test_project/main.tf')
    repo.git.commit('-m', 'Version 1')
    first_commit = repo.head.commit.hexsha
    
    # 두 번째 버전
    with open(test_file_path, 'w') as f:
        f.write('Version 2 content')
    repo.git.add('test_project/main.tf')
    repo.git.commit('-m', 'Version 2')
    
    # 백업 옵션으로 롤백
    result = runner.invoke(cli, ['rollback', 'test_project', first_commit, '--backup'], input='n\n')
    assert result.exit_code == 0
    assert "백업 브랜치" in result.output
    
    # 백업 브랜치 존재 확인
    backup_branches = [ref.name for ref in repo.refs if ref.name.startswith('backup_')]
    assert len(backup_branches) == 1

def test_rollback_with_local_changes(mock_home_dir, mock_repo):
    """로컬 변경사항이 있는 경우의 롤백 테스트"""
    repo, repo_path = mock_repo
    runner = CliRunner()
    
    # 초기 설정
    runner.invoke(cli, ['init', repo_path, 'test_project', 'terraform'])
    
    workspace = os.path.join(mock_home_dir, "workspace")
    project_path = os.path.join(workspace, "test_project")
    test_file_path = os.path.join(project_path, "test_project", "main.tf")
    
    repo = git.Repo(project_path)
    
    # 첫 번째 커밋
    with open(test_file_path, 'w') as f:
        f.write('Initial content')
    repo.git.add('test_project/main.tf')
    repo.git.commit('-m', 'Initial commit')
    first_commit = repo.head.commit.hexsha
    
    # 로컬 변경사항 만들기
    with open(test_file_path, 'w') as f:
        f.write('Modified content')
    
    # 롤백 시도
    result = runner.invoke(cli, ['rollback', 'test_project', first_commit])
    assert result.exit_code == 0
    assert "저장되지 않은 로컬 변경사항이 있습니다" in result.output
    
    # force 옵션으로 롤백
    result = runner.invoke(cli, ['rollback', 'test_project', first_commit, '--force'], input='n\n')
    assert result.exit_code == 0
    assert f"버전 {first_commit}으로 롤백했습니다" in result.output

def test_rollback_nonexistent_version(mock_home_dir, mock_repo):
    """존재하지 않는 버전으로의 롤백 테스트"""
    repo, repo_path = mock_repo
    runner = CliRunner()
    
    # 초기 설정
    runner.invoke(cli, ['init', repo_path, 'test_project', 'terraform'])
    
    # 존재하지 않는 버전으로 롤백 시도
    result = runner.invoke(cli, ['rollback', 'test_project', 'nonexistent'])
    assert result.exit_code == 0
    assert "버전 nonexistent을 찾을 수 없습니다" in result.output
