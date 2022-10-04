import { Navigate, useRoutes } from 'react-router-dom';
// layouts
import DashboardLayout from './layouts/dashboard';
import LogoOnlyLayout from './layouts/LogoOnlyLayout';
//
import Resources from './pages/Resources';
import NotFound from './pages/Page404';
import DashboardApp from './pages/DashboardApp';
import ComingSoon from './pages/PageComingSoon';

// ----------------------------------------------------------------------

export default function Router() {
  return useRoutes([
    {
      path: '/dashboard',
      element: <DashboardLayout />,
      children: [
        { 
        path: 'warm/', 
        children:[
          { path: 'aws', element: <DashboardApp experimentType={'warm'}/> },
          { path: 'google', element: <ComingSoon /> }
        ]},
        { 
        path: 'cold', 
        children:[
          { path: 'baseline', element: <ComingSoon /> },
          { path: 'image-size', element: <ComingSoon /> },
          { path: 'deployment-language', element: <ComingSoon /> }
        ]},
        { 
        path: 'data-transfer', 
        children:[
          { path: 'inline', element: <ComingSoon /> },
          { path: 'storage-based', element: <ComingSoon /> },
        ] 
        },
        {
        path: 'bursty', 
        children:[
          { path: 'short-iat', element: <ComingSoon /> },
          { path: 'long-iat', element: <ComingSoon /> },
          { path: 'scheduling-policy', element: <ComingSoon /> },
        ] 
      },
        { path: 'resources', element: <Resources /> },
      ],

    },
    {
      path: '/',
      element: <LogoOnlyLayout />,
      children: [
        { path: '/', element: <Navigate to="/dashboard/warm/aws" /> },
        { path: '/dashboard', element: <Navigate to="/dashboard/warm/aws" /> },
        { path: '404', element: <NotFound /> },
        { path: '*', element: <Navigate to="/404" /> },
      ],
    },
    { path: '*', element: <Navigate to="/404" replace /> },
  ]);
}
