import click
import os
import git
from datetime import datetime
from ..utils.system import get_workspace_dir

@click.command()
@click.argument('project_name')
@click.option('--limit', default=10, help='조회할 히스토리 개수')
@click.option('--format', type=click.Choice(['simple', 'detailed']), default='simple', help='출력 형식')
def history(project_name, limit, format):
    """프로젝트의 버전 히스토리를 조회"""
    try:
        workspace = get_workspace_dir()
        project_path = os.path.join(workspace, project_name)
        
        # 프로젝트 존재 여부 확인
        if not os.path.exists(project_path):
            print(f"프로젝트 {project_name}을 찾을 수 없습니다.")
            return False
            
        repo = git.Repo(project_path)
        
        # 프로젝트 브랜치의 커밋 히스토리 조회
        commits = list(repo.iter_commits(project_name, max_count=limit))
        
        if not commits:
            print(f"{project_name} 프로젝트의 커밋 히스토리가 없습니다.")
            return True
            
        print(f"\n{project_name} 프로젝트 버전 히스토리:")
        print("=" * 50)
        
        for commit in commits:
            date = datetime.fromtimestamp(commit.committed_date)
            
            if format == 'simple':
                print(f"커밋: {commit.hexsha[:8]}")
                print(f"날짜: {date.strftime('%Y-%m-%d %H:%M:%S')}")
                print(f"메시지: {commit.message.strip()}")
                print("-" * 50)
            else:
                print(f"커밋: {commit.hexsha}")
                print(f"작성자: {commit.author.name} <{commit.author.email}>")
                print(f"날짜: {date.strftime('%Y-%m-%d %H:%M:%S')}")
                print(f"메시지:\n{commit.message.strip()}")
                
                # 변경된 파일 목록
                files = commit.stats.files
                if files:
                    print("\n변경된 파일:")
                    for file_path, stats in files.items():
                        print(f"- {file_path} (+{stats['insertions']}, -{stats['deletions']})")
                print("=" * 50)
        
        return True
        
    except git.exc.GitCommandError as e:
        print(f"Git 작업 중 오류 발생: {e}")
        return False
    except Exception as e:
        print(f"히스토리 조회 실패: {e}")
        return False
