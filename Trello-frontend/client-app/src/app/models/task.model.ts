export interface Task {
    id: string;
    projectId: string;
    name: string;
    description: string;
    assignedTo: string; // `primitive.ObjectID` je zamijenjen sa `string`
    createdAt: Date;
    dueDate: Date;
    status: string;
    priority: string;
    documentLinks: string[];
    workflowStateId: string;
  }