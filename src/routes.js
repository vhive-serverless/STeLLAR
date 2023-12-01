import { Navigate, useRoutes } from 'react-router-dom';
// layouts
import DashboardLayout from './layouts/dashboard';
import LogoOnlyLayout from './layouts/LogoOnlyLayout';
//
import Team from './pages/Team';
import NotFound from './pages/Page404';
import BaselineLatencyWarm from './pages/BaselineLatencyWarm';
import BaselineLatencyCold from './pages/BaselineLatencyColdNew';
import BaselineLatencyColdImageSize from './pages/BaselineLatencyColdImageSizeAll';
// import BaselineLatencyColdDeploymentLanguage from './pages/BaselineLatencyDeploymentLanguage';
import ComingSoon from './pages/PageComingSoon';
import About from './pages/About';

// ----------------------------------------------------------------------

export default function Router() {
  return useRoutes([
    {
      path: '/dashboard',
      element: <DashboardLayout />,
      children: [
        { 
          path: 'about/', 
   element: <About/> 
  },
        { 
        path: 'warm/', 
        children:[
          { path: 'aws', element: <BaselineLatencyWarm experimentType={'warm-baseline-aws'}/> },
        ]},
        { 
        path: 'cold', 
        children:[
          { path: 'baseline', element: <BaselineLatencyCold/> },
          // { path: 'image-size', element:<BaselineLatencyColdImageSize/> },
          { path: 'image-size', element:<BaselineLatencyColdImageSize/> },
          // { path: 'deployment-language', element: <BaselineLatencyColdDeploymentLanguage /> }
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
        { path: 'team', element: <Team /> },
      ],

    },
    {
      path: '/',
      element: <LogoOnlyLayout />,
      children: [
        { path: '/', element: <Navigate to="/dashboard/about" /> },
        { path: '/dashboard', element: <Navigate to="/dashboard/about" /> },
        { path: '404', element: <NotFound /> },
        { path: '*', element: <Navigate to="/404" /> },
      ],
    },
    { path: '*', element: <Navigate to="/404" replace /> },
  ]);
}
