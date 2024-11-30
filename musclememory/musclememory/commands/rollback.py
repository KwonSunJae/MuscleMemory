import click
import os
import git
from ..utils.system import get_workspace_dir
from ..models.project import Project

@click.command()
@click.argument('project_name')
@click.argument('version')
@click.option('--force', is_flag=True, help='강제 롤백 여부')
@click.option('--backup', is_flag=True, help='현재 상태 백업 여부')
def rollback(project_name, version, force, backup):
    """특정 버전으로 롤백"""
    try:
        workspace = get_workspace_dir()
        project_path = os.path.join(workspace, project_name)
        
        # 프로젝트 존재 여부 확인
        if not os.path.exists(project_path):
            print(f"프로젝트 {project_name}을 찾을 수 없습니다.")
            return False
            
        repo = git.Repo(project_path)
        
        # 현재 브랜치 확인
        current = repo.active_branch.name
        if current != project_name:
            print(f"프로젝트 브랜치({project_name})로 전환합니다.")
            repo.git.checkout(project_name)
        
        # 변경사항 확인
        if repo.is_dirty() and not force:
            print("저장되지 않은 로컬 변경사항이 있습니다. --force 옵션을 사용하거나 변경사항을 처리해주세요.")
            return False
            
        # 백업 브랜치 생성
        if backup:
            backup_branch = f"backup_{project_name}_{repo.head.commit.hexsha[:8]}"
            repo.git.branch(backup_branch)
            print(f"현재 상태를 {backup_branch} 브랜치에 백업했습니다.")
        
        try:
            # 특정 버전으로 롤백
            repo.git.reset('--hard', version)
            
            # 원격 저장소에 강제 푸시
            if click.confirm('변경사항을 원격 저장소에 푸시하시겠습니까?'):
                repo.git.push('origin', project_name, '--force')
                print(f"변경사항을 원격 저장소에 푸시했습니다.")
            
            print(f"{project_name}을 버전 {version}으로 롤백했습니다.")
            return True
            
        except git.exc.GitCommandError:
            print(f"버전 {version}을 찾을 수 없습니다.")
            if backup:
                # 백업 브랜치로 복구
                repo.git.reset('--hard', backup_branch)
                print(f"백업 브랜치({backup_branch})로 복구했습니다.")
            return False
            
    except git.exc.GitCommandError as e:
        print(f"Git 작업 중 오류 발생: {e}")
        return False
    except Exception as e:
        print(f"롤백 실패: {e}")
        return False
