import pytest
import os
import shutil
import tempfile
import git
from pathlib import Path

@pytest.fixture
def temp_dir():
    """임시 디렉토리 생성"""
    temp_path = tempfile.mkdtemp()
    yield temp_path
    shutil.rmtree(temp_path)

@pytest.fixture
def mock_home_dir(temp_dir, monkeypatch):
    """임시 홈 디렉토리 설정"""
    muscle_dir = os.path.join(temp_dir, ".muscle")
    os.makedirs(muscle_dir)
    monkeypatch.setenv("HOME", temp_dir)
    return muscle_dir

@pytest.fixture
def mock_repo(temp_dir):
    """테스트용 Git 저장소 생성"""
    repo_path = os.path.join(temp_dir, "test_repo")
    repo = git.Repo.init(repo_path)
    return repo, repo_path
