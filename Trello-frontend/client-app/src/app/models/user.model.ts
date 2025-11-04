export type UserRole = 'manager' | 'member' | 'unauthenticated';

export interface User {
  id: string;
  username: string;
  first_name: string;
  last_name: string;
  email: string;
//   password?: string;
  role: UserRole;
  createdAt: Date;
  updatedAt: Date;
  isActive: boolean;
  code: string;
}
