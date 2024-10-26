import os
import json
from datetime import datetime, timedelta

class Lock:
    def __init__(self, project, owner):
        self.project = project
        self.owner = owner
        self.lock_path = os.path.join(project.inner_path, ".lock")
        
    def acquire(self, timeout_minutes=30):
        """잠금 획득"""
        expired_time = (datetime.now() + timedelta(minutes=timeout_minutes)).isoformat()
        lock_data = {
            'owner': self.owner,
            'expired_time': expired_time
        }
        with open(self.lock_path, 'w') as f:
            json.dump(lock_data, f)
            
    def release(self):
        """잠금 해제"""
        if os.path.exists(self.lock_path):
            os.remove(self.lock_path)
            
    @property
    def exists(self):
        return os.path.exists(self.lock_path)
        
    def get_info(self):
        """잠금 정보 조회"""
        if not self.exists:
            return None
        with open(self.lock_path, 'r') as f:
            return json.load(f)
