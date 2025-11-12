import { lazy, Suspense } from "react";
import { Atom }from "react-loading-indicators"
import { Outlet, useRoutes } from "react-router-dom"
import DashboardLayout from "../layouts/Dashboard/Dashboard";

export const Homepage = lazy(()=>import("../pages/Home/Home"))

export default function Router() {
  const loading = <Atom color="#F90627" size="large"/>

  return useRoutes([
    {
      element: (
        <DashboardLayout>
          <Suspense fallback={loading}>
            <Outlet/>
          </Suspense>
        </DashboardLayout>
      ),
      children: [
        { element: <Homepage/>, index: true },
      ]
    }
  ]);
}
