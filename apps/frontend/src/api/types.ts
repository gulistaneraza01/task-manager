import { apiContract } from "@task-manager/openapi/contracts";
import type { ServerInferRequest } from "@ts-rest/core";

export type TRequests = ServerInferRequest<typeof apiContract>;
