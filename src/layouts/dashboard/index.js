import { useState } from 'react';
import { Outlet } from 'react-router-dom';
// material
import{ IconButton } from '@mui/material';
import { styled } from '@mui/material/styles';
import DashboardSidebar from './DashboardSidebar';


import Iconify from '../../components/Iconify';

// ----------------------------------------------------------------------

const APP_BAR_MOBILE = 64;
const APP_BAR_DESKTOP = 20;

const RootStyle = styled('div')({
  display: 'flex',
  minHeight: '100%',
  overflow: 'hidden'
});

const MainStyle = styled('div')(({ theme }) => ({
  flexGrow: 1,
  overflow: 'auto',
  minHeight: '100%',
  paddingTop: APP_BAR_MOBILE + 24,
  paddingBottom: theme.spacing(10),
  [theme.breakpoints.up('lg')]: {
    paddingTop: APP_BAR_DESKTOP + 24,
    paddingLeft: theme.spacing(2),
    paddingRight: theme.spacing(2)
  }
}));

// ----------------------------------------------------------------------

export default function DashboardLayout() {
  const [open, setOpen] = useState(false);

  return (<>
    <IconButton onClick={()=>setOpen(true)} sx={{ mt: 1, color: 'text.primary',display: { lg: 'none' } }}>
          <Iconify icon="eva:menu-2-fill" />
        </IconButton>
    <RootStyle>
      <DashboardSidebar isOpenSidebar={open} onCloseSidebar={() => setOpen(false)} />
      <MainStyle>
        <Outlet />
      </MainStyle>
    </RootStyle>
    </>
  );
}
