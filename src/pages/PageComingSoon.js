import { Link as RouterLink } from 'react-router-dom';
// @mui
import { styled } from '@mui/material/styles';
import { Button, Typography, Container, Box } from '@mui/material';
// components
import Page from '../components/Page';

// ----------------------------------------------------------------------

const ContentStyle = styled('div')(() => ({
  maxWidth: 480,
  margin: 'auto',
  minHeight: '90vh',
  display: 'flex',
  justifyContent: 'center',
  flexDirection: 'column',
}));

// ----------------------------------------------------------------------

export default function PageComingSoon() {
  return (
    <Page title="Statistics Coming Soon!">
      <Container>
        <ContentStyle sx={{ textAlign: 'center', alignItems: 'center' }}>
          <Typography variant="h3" paragraph>
              Coming Soon!
          </Typography>

          <Typography sx={{ color: 'text.secondary' }}>
            We are constantly working on to make these results available to you as soon as possible.
            We'll notify you once we are ready :)
          </Typography>

          <Box
            component="img"
            src="/STeLLAR/dynamic/graph.gif"
            sx={{ height: 300, mx: 'auto', my: { xs: 5, sm: 10 } }}
          />

          <Button to="/" size="large" variant="contained" component={RouterLink}>
            Go to Home
          </Button>
        </ContentStyle>
      </Container>
    </Page>
  );
}
