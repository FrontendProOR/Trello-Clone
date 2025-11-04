export interface Document {
    id: string;
    projectId: string;
    taskId?: string;
    name: string;
    url: string;
    uploadedAt: Date;
    uploadedBy: string;
  }