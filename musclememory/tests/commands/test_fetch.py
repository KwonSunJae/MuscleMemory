import pytest
import os
import git
from click.testing import CliRunner
from musclememory.cli import cli

def test_fetch_latest(mock_home_dir, mock_repo):
    """최신 IaC 파일 불러오기 테스트"""
    repo, repo_path = mock_repo
    runner = CliRunner()
    
    # 초기 설정
    runner.invoke(cli, ['init', repo_path, 'test_project', 'terraform'])
    
    # 테스트 파일 생성 및 커밋
    workspace = os.path.join(mock_home_dir, "workspace")
    project_path = os.path.join(workspace, "test_project")
    test_file_path = os.path.join(project_path, "test_project", "main.tf")
    
    with open(test_file_path, 'w') as f:
        f.write('resource "aws_instance" "example" {}')
    
    repo = git.Repo(project_path)
    repo.git.add('test_project/main.tf')
    repo.git.commit('-m', 'Add test file')
    
    # fetch latest 테스트
    result = runner.invoke(cli, ['fetch', 'latest', 'test_project'])
    assert result.exit_code == 0
    assert "최신 IaC 파일을 가져왔습니다" in result.output

def test_fetch_version(mock_home_dir, mock_repo):
    """특정 버전 IaC 파일 불러오기 테스트"""
    repo, repo_path = mock_repo
    runner = CliRunner()
    
    # 초기 설정
    runner.invoke(cli, ['init', repo_path, 'test_project', 'terraform'])
    
    # 여러 버전의 파일 생성 및 커밋
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
    
    # 특정 버전 가져오기 테스트
    result = runner.invoke(cli, ['fetch', 'version', 'test_project', first_commit])
    assert result.exit_code == 0
    assert f"{first_commit} 버전 IaC 파일을 가져왔습니다" in result.output
    
    # 파일 내용 확인
    with open(test_file_path, 'r') as f:
        content = f.read()
    assert 't2.micro' in content

def test_fetch_nonexistent_version(mock_home_dir, mock_repo):
    """존재하지 않는 버전 테스트"""
    repo, repo_path = mock_repo
    runner = CliRunner()
    
    # 초기 설정
    runner.invoke(cli, ['init', repo_path, 'test_project', 'terraform'])
    
    # 존재하지 않는 버전 테스트
    result = runner.invoke(cli, ['fetch', 'version', 'test_project', 'nonexistent'])
    assert result.exit_code == 0
    assert "버전 nonexistent을 찾을 수 없습니다" in result.output

def test_fetch_with_local_changes(mock_home_dir, mock_repo):
    """로컬 변경사항이 있는 경우 테스트"""
    repo, repo_path = mock_repo
    runner = CliRunner()
    
    # 초기 설정
    runner.invoke(cli, ['init', repo_path, 'test_project', 'terraform'])
    
    # 테스트 파일 생성 및 변경
    workspace = os.path.join(mock_home_dir, "workspace")
    project_path = os.path.join(workspace, "test_project")
    test_file_path = os.path.join(project_path, "test_project", "main.tf")
    
    with open(test_file_path, 'w') as f:
        f.write('resource "aws_instance" "example" {}')
    
    # fetch 시도
    result = runner.invoke(cli, ['fetch', 'latest', 'test_project'])
    assert result.exit_code == 0
    assert "저장되지 않은 로컬 변경사항이 있습니다" in result.output
