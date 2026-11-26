import { z } from "zod";

export const ZTodoComment = z.object({
  id: z.string().uuid(),
  todoId: z.string().uuid(),
  userId: z.string(),
  content: z.string(),
  parentCommentId: z.string().uuid().nullable(),
  createdAt: z.string(),
  updatedAt: z.string(),
});
