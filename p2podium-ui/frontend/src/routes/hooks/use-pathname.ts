import { useNavigate, useLocation, useParams } from "react-router-dom";

export function useRouter() {
  const navigate = useNavigate();
  return {
    push: (path: string) => navigate(path),
    back: () => navigate(-1),
  };
}

export function useRouteParams<T>() {
  return useParams() as T;
}

export function usePathname() {
  const { pathname } = useLocation();
  return pathname;
}
