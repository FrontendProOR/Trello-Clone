export interface Workflow {
    id: string;
    projectId: string;
    graphData: string; // JSON enkodirani graf
    createdAt: Date;
    lastUpdated: Date;
  }