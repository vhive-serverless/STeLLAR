// scroll bar
import 'simplebar/src/simplebar.css';

import ReactDOM from 'react-dom';
import { HashRouter } from 'react-router-dom';
import { HelmetProvider } from 'react-helmet-async';
import ReactGA4 from 'react-ga4';

//
import App from './App';
import * as serviceWorker from './serviceWorker';
import reportWebVitals from './reportWebVitals';

// ----------------------------------------------------------------------

const TRACKING_ID = "G-56TWKV1598"; // OUR_TRACKING_ID
ReactGA4.initialize(TRACKING_ID);
ReactDOM.render(
  <HelmetProvider>
    <HashRouter>
      <App />
    </HashRouter>
  </HelmetProvider>,
  document.getElementById('root')
);

// If you want to enable client cache, register instead.
serviceWorker.unregister();

// If you want to start measuring performance in your app, pass a function
// to log results (for example: reportWebVitals(console.log))
// or send to an analytics endpoint. Learn more: https://bit.ly/CRA-vitals
reportWebVitals();
