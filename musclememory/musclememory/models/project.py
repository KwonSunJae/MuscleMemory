import os
from ..utils.system import get_workspace_dir
from ..utils.git import init_repo, ensure_branch

class Project:
    def __init__(self, name, repo_url=None):
        self.name = name
        self.repo_url = repo_url
        self.workspace = get_workspace_dir()
        self.project_path = os.path.join(self.workspace, name)
        self.inner_path = os.path.join(self.project_path, name)
        
    @property
    def exists(self):
        return os.path.exists(self.project_path)
        
    def init_repository(self):
        """저장소 초기화"""
        repo, is_new = init_repo(self.repo_url, self.project_path)
        return repo, is_new
        
    def ensure_directories(self):
        """프로젝트 디렉토리 구조 생성"""
        os.makedirs(self.inner_path, exist_ok=True)
        
    def get_repo(self):
        """Git 저장소 객체 반환"""
        import git
        return git.Repo(self.project_path)
        
    def ensure_branch(self, repo, branch_name, upstream=True):
        """브랜치 생성 또는 체크아웃"""
        try:
            branch = getattr(repo.heads, branch_name)
            branch.checkout()
            return branch, False
        except AttributeError:
            branch = repo.create_head(branch_name)
            branch.checkout()
            if upstream:
                repo.git.push('--set-upstream', 'origin', branch_name)
            return branch, True