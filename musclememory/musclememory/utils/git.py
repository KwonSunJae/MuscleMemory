import git
import os

def init_repo(repo_url, project_path):
    """저장소 초기화 또는 가져오기"""
    if not os.path.exists(project_path):
        repo = git.Repo.clone_from(repo_url, project_path)
        is_new_clone = True
    else:
        repo = git.Repo(project_path)
        is_new_clone = False
    return repo, is_new_clone

def ensure_branch(repo, branch_name, upstream=True):
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
