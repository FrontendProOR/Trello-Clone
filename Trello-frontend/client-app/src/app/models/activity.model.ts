export interface Activity {
    id: string;
    userId: string;
    projectId: string;
    taskId?: string; 
    action: string;
    timestamp: Date;
    description: string;
  }