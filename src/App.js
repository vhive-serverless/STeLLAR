import { LocalizationProvider } from '@mui/x-date-pickers';
import { AdapterDateFns } from '@mui/x-date-pickers/AdapterDateFns';
// routes
import Router from './routes';
// theme
import ThemeProvider from './theme';
// components
import ScrollToTop from './components/ScrollToTop';
import { BaseOptionChartStyle } from './components/chart/BaseOptionChart';

// ----------------------------------------------------------------------

export default function App() {
  return (
    <ThemeProvider>
        <LocalizationProvider dateAdapter={AdapterDateFns}>
      <ScrollToTop />
      <BaseOptionChartStyle />
      <Router />
        </LocalizationProvider>
    </ThemeProvider>
  );
}
