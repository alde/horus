import React from "react";
import PropTypes from "prop-types";
import { withStyles } from "@material-ui/core/styles";
import AppBar from "@material-ui/core/AppBar";
import Toolbar from "@material-ui/core/Toolbar";
import Typography from "@material-ui/core/Typography";
import IconButton from "@material-ui/core/IconButton";
import { Paper } from "@material-ui/core";

const styles = {
  root: {
    flexGrow: 1
  },
  menuButton: {
    marginLeft: -18,
    marginRight: 10
  }
};

function App(props) {
  const { classes } = props;
  return (
    <div className={classes.root}>
      <AppBar position="static">
        <Toolbar variant="dense">
          <IconButton
            className={classes.menuButton}
            color="inherit"
            aria-label="Menu"
          >
            <img alt="horus logo" src="/eye_of_horus_small.png" />
          </IconButton>
          <Typography variant="h6" color="inherit">
            Horus
          </Typography>
        </Toolbar>
      </AppBar>
      <Paper>
        <Typography>There should be some sort of UI</Typography>
      </Paper>
    </div>
  );
}

App.propTypes = {
  classes: PropTypes.object.isRequired
};

export default withStyles(styles)(App);
