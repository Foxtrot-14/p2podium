import { useParams } from "react-router-dom";

export function useParamsValue<T extends Record<string, string>>() {
  return useParams<T>();
}
