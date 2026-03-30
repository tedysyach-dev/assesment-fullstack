export type ApiResponse<T> = {
  status: boolean;
  message: string;
  resource: T;
};
