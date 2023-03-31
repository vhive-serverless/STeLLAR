import { useState } from 'react';
// material
import { Grid, Card,Box, Container, Stack,TextField, Typography,Link,Fab,Popover,Button } from '@mui/material';
// components
import Page from '../components/Page';
import AccountProfile from '../components/AccountProfile'
// data
import { members } from '../data/members';
// ----------------------------------------------------------------------

const styles = {
  container: {
    display: 'flex',
    justifyContent: 'center',
    alignItems:'center'
  
  },
  card: {
    width:'100%',
  },
  fab:{
    position:'absolute',
    botton:0,
  }
};

export default function Info() {

  const [anchorEl, setAnchorEl] = useState(null);
  const [feedback, setFeedback] = useState('');
  const [email, setEmail] = useState('');

  const handleClick = (event) => {
    setAnchorEl(event.currentTarget);
  };

  const handleClose = () => {
    setAnchorEl(null);
  };

  const handleSubmit = (event) => {
    event.preventDefault();
    // Submit feedback and email
    console.log(`Feedback: ${feedback}, Email: ${email}`);
    // Clear form fields
    setFeedback('');
    setEmail('');
    handleClose();
  };

  const open = Boolean(anchorEl);
  const id = open ? 'popup' : undefined;

  return (
    <Page title="Dashboard: Info" sx={styles.container}>
      <Container maxWidth="xl">
        
        <Stack direction="row" alignItems="center" justifyContent="space-between" mt={3}>
          <Typography variant="h4" gutterBottom>
            Team
          </Typography>
        </Stack>
          <Grid container xs={12}>
                {members.map((member) => (
                   <Grid item lg={3}
                   md={6}
                   xs={12}
                   sx={{pr:2,pb:2}} key={member.name} >
                    <Card>
            <AccountProfile member={member}/>
            </Card>
            </Grid>
            ))}

            </Grid>
        <Card sx={styles.card} > 

          <Stack direction="row" alignItems="center" justifyContent="center" sx={{paddingY:0.9}}>                 
          <Typography
            color="textPrimary"
            sx={{fontSize:16,fontWeight:500,marginX:2}}
            align='center'
          >
             Find us on :
          </Typography>     
            <Link target="_blank" href={'https://github.com/vhive-serverless'}><Box component="img" src="/STeLLAR/static/icons/github.png" sx={{height: '30px',marginX:2 }}/></Link>
            <Link target="_blank" href={'https://www.youtube.com/playlist?list=PLVdxPJaekjWqBsEUwnrYRQCaMqvcDVsBE'}><Box component="img" src="/STeLLAR/static/icons/youtube.png" sx={{height: '30px',marginX:2  }}/></Link>
            <Link target="_blank" href={'https://join.slack.com/t/vhivetutorials/shared_invite/zt-1fk4v71gn-nV5oev5sc9F4fePg3_OZMQ'}><Box component="img" src="/STeLLAR/static/icons/slack.png" sx={{height: '30px',marginX:2  }}/></Link>
          </Stack>

            </Card>
      </Container>
      <div style={{ position: 'absolute', bottom: 10, right: 10 }}>
      <Fab color="primary" sx={{ borderRadius: '16px',width: '220px',height:'20px' }}  onClick={handleClick}>
      <span style={{ display: 'flex', justifyContent: 'center', alignItems: 'center' }}>
      <Typography
            color="white"
            sx={{fontSize:16,fontWeight:500,marginX:2}}
            align='center'
          >
             Got some feedback ?
        </Typography>    
          </span>
      </Fab>
      <Popover
        id={id}
        open={open}
        anchorEl={anchorEl}
        onClose={handleClose}
        anchorOrigin={{
          vertical: 'bottom',
          horizontal: 'center',
        }}
        transformOrigin={{
          vertical: 'top',
          horizontal: 'center',
        }}
      >
         <Box sx={{padding:2}}>
       <form onSubmit={handleSubmit}>
       <Typography
            color="primary"
            sx={{fontSize:16,fontWeight:600}}
            align="center"
          >
             Help us to improve 
        </Typography>   
       <TextField
            label="Email"
            type="email"
            value={email}
            onChange={(event) => setEmail(event.target.value)}
            variant="outlined"
            margin="normal"
            fullWidth
            required
          />
          <TextField
            label="Feedback"
            multiline
            value={feedback}
            onChange={(event) => setFeedback(event.target.value)}
            variant="outlined"
            margin="normal"
            fullWidth
            required
          />
          
          <Button type="submit" variant="contained" color="primary" sx={{marginTop:'5px'}}>
            Submit
          </Button>
        </form>
        </Box>
      </Popover>
    </div>
    </Page>
  
  );
}
