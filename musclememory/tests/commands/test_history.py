import pytest
import os
import git
from click.testing import CliRunner
from musclememory.cli import cli

def test_history_basic(mock_home_dir, mock_repo):
    """기본 히스토리 조회 테스트"""
    repo, repo_path = mock_repo
    runner = CliRunner()
    
    # 초기 설정
    runner.invoke(cli, ['init', repo_path, 'test_project', 'terraform'])
    
    # 여러 커밋 생성
    workspace = os.path.join(mock_home_dir, "workspace")
    project_path = os.path.join(workspace, "test_project")
    test_file_path = os.path.join(project_path, "test_project", "main.tf")
    
    repo = git.Repo(project_path)
    
    # 첫 번째 커밋
    with open(test_file_path, 'w') as f:
        f.write('resource "aws_instance" "example" { instance_type = "t2.micro" }')
    repo.git.add('test_project/main.tf')
    repo.git.commit('-m', 'First commit')
    
    # 두 번째 커밋
    with open(test_file_path, 'w') as f:
        f.write('resource "aws_instance" "example" { instance_type = "t2.small" }')
    repo.git.add('test_project/main.tf')
    repo.git.commit('-m', 'Second commit')
    
    # 히스토리 조회 테스트
    result = runner.invoke(cli, ['history', 'test_project'])
    assert result.exit_code == 0
    assert 'First commit' in result.output
    assert 'Second commit' in result.output

def test_history_with_limit(mock_home_dir, mock_repo):
    """히스토리 개수 제한 테스트"""
    repo, repo_path = mock_repo
    runner = CliRunner()
    
    # 초기 설정 및 커밋 생성
    runner.invoke(cli, ['init', repo_path, 'test_project', 'terraform'])
    workspace = os.path.join(mock_home_dir, "workspace")
    project_path = os.path.join(workspace, "test_project")
    test_file_path = os.path.join(project_path, "test_project", "main.tf")
    
    repo = git.Repo(project_path)
    
    # 3개의 커밋 생성
    for i in range(3):
        with open(test_file_path, 'w') as f:
            f.write(f'resource "aws_instance" "example_{i}" {{}}')
        repo.git.add('test_project/main.tf')
        repo.git.commit('-m', f'Commit {i+1}')
    
    # 2개로 제한하여 조회
    result = runner.invoke(cli, ['history', 'test_project', '--limit', '2'])
    assert result.exit_code == 0
    assert 'Commit 3' in result.output
    assert 'Commit 2' in result.output
    assert 'Commit 1' not in result.output

def test_history_detailed_format(mock_home_dir, mock_repo):
    """상세 형식 히스토리 조회 테스트"""
    repo, repo_path = mock_repo
    runner = CliRunner()
    
    # 초기 설정
    runner.invoke(cli, ['init', repo_path, 'test_project', 'terraform'])
    workspace = os.path.join(mock_home_dir, "workspace")
    project_path = os.path.join(workspace, "test_project")
    test_file_path = os.path.join(project_path, "test_project", "main.tf")
    
    repo = git.Repo(project_path)
    
    # 커밋 생성
    with open(test_file_path, 'w') as f:
        f.write('resource "aws_instance" "example" {}')
    repo.git.add('test_project/main.tf')
    repo.git.commit('-m', 'Test commit')
    
    # 상세 형식으로 조회
    result = runner.invoke(cli, ['history', 'test_project', '--format', 'detailed'])
    assert result.exit_code == 0
    assert 'Test commit' in result.output
    assert 'main.tf' in result.output
    assert '변경된 파일:' in result.output
