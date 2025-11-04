import { Task } from './task.model';
import { Document } from './document.model';
import { Activity } from './activity.model';
import { ProjectAnalytics } from './project-analytics.model';

export interface Project {
  id: string;
  name: string;
  description: string;
  managerId: string;
  createdAt: Date;
  expectedEndDate: Date;
  minMembers: number;
  maxMembers: number;
  status: string;
  members: string[]; // Lista korisničkih ID-ova
  tasks: Task[];
  workflowGraphId: string;
  documents: Document[];
  notificationsEnabled: boolean;
  history: Activity[];
  analytics: ProjectAnalytics;
}