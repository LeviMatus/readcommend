import React from "react";
import ReactDOM from "react-dom";
import CssBaseline from "@material-ui/core/CssBaseline";
import { ThemeProvider } from "@material-ui/core/styles";
import App from "./components/App";
import Banner from "./components/Banner";
import theme from "./theme";

ReactDOM.render(
  <div>
    <Banner />
    <ThemeProvider theme={theme}>
      <CssBaseline />
      <App />
    </ThemeProvider>
  </div>,
  document.querySelector("#root")
);
