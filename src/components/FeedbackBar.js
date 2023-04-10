import {useState } from 'react';
import axios from 'axios';
// material
import { Box ,Alert,TextField,Snackbar, Typography,Fab,Popover,Button } from '@mui/material';
// components


export default function FeedbackBar() {
    const emailAPI =  'https://o67aj2z412.execute-api.us-west-2.amazonaws.com/default/send-email-stellar';
    


    const [feedback, setFeedback] = useState('');
    const [email, setEmail] = useState('');
    const [name, setName] = useState('');
    const [openSnackbarSuccess, setOpenSnackbarSuccess] = useState(false);
    const [openSnackbarError, setOpenSnackbarError] = useState(false);
    const [anchorEl, setAnchorEl] = useState(null);

    const open = Boolean(anchorEl);
    const id = open ? 'popup' : undefined;
    
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
    
    axios.post(emailAPI, {
        name,email,feedback
    })
        .then(response => {
        setOpenSnackbarSuccess(true);
        console.log(response);
        })
        .catch(error => {
        setOpenSnackbarError(true);
        console.log(error);
        });

    // Clear form fields
    setFeedback('');
    setEmail('');
    setName('');
    handleClose();
    };
        
    const handleCloseSnackbar = () => {
    setOpenSnackbarSuccess(false);
    setOpenSnackbarError(false);
    };
          
    return (
    <>
    <div style={{ position: 'relative' }}>
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
        sx={{backgroundColor: 'rgba(0, 0, 0, 0.5)',transition: 'opacity 0.5s ease-in-out' }}
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
            label="Name"
            value={name}
            onChange={(event) => setName(event.target.value)}
            variant="outlined"
            margin="normal"
            fullWidth
            required
          />   
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
      {openSnackbarSuccess && <Snackbar open={openSnackbarSuccess} autoHideDuration={6000} onClose={handleCloseSnackbar}>
         <Alert elevation={6} variant="filled" onClose={handleCloseSnackbar} severity="success">
          Feedback sent successfully!
        </Alert>
      </Snackbar>
      }
      {openSnackbarError && <Snackbar open={openSnackbarError} autoHideDuration={6000} onClose={handleCloseSnackbar}>
         <Alert elevation={6} variant="filled" onClose={handleCloseSnackbar} severity="error">
          An error occured!
        </Alert>
      </Snackbar>
      }
    </div>
    </>);
};