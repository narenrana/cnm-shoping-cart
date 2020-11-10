import React from "react";
import CssBaseline from "@material-ui/core/CssBaseline";
import Grid from "@material-ui/core/Grid";
import Container from "@material-ui/core/Container";
import Header from "./header/Header";
import Footer from "./footer/Footer";
import { ThemeProvider } from "@material-ui/core/styles";
import SideNavBar from "./sidebar/SideNavBar";
import ProductList from "./products/components/product-home/ProductHome";
import { Checkout } from "./checkout";
import { SignIn } from "./sign-in";
import theme from "./theme";
import {
  BrowserRouter as Router,
  Switch,
  Route,
  Redirect,
  useHistory,
} from "react-router-dom";

import { useSelector } from "react-redux";

const sections = [
  { title: "SignIn", url: "/login" },
  { title: "Home", url: "/home" },
  { title: "Orders", url: "/orders" },
  { title: "Paymemnts", url: "/paymemnts" },
];

export default function App() {
  let historyq = useHistory();
  console.log({ historyq });
  const { auth = {} } = useSelector((state) => state.commonStore);
  return (
    <ThemeProvider theme={theme}>
      <CssBaseline />{" "}
      <Router>
        <Container maxWidth="lg">
          {auth.isLogin && <Header title="Cart" sections={sections} />}

          <Switch>
            <Route path="/login">
              <SignIn />
            </Route>
            <PrivateRoute path="/payment">
              <div>Payment</div>
            </PrivateRoute>
            <PrivateRoute path="/checkout">
              <main>
                <Checkout />
              </main>
            </PrivateRoute>
            <PrivateRoute path="/">
              <main>
                <Grid container>
                  <Grid item xs={12} md={3} xl={3}>
                    <SideNavBar />
                  </Grid>
                  <Grid item xs={12} md={9} xl={9}>
                    <ProductList />
                  </Grid>
                </Grid>
              </main>
            </PrivateRoute>
          </Switch>
        </Container>
        {auth.isLogin && (
          <Footer
            title="Footer"
            description="Something here to give the footer a purpose!"
          />
        )}
      </Router>
    </ThemeProvider>
  );
}

function PrivateRoute({ children, ...rest }) {
  const { auth = {} } = useSelector((state) => state.commonStore);

  return (
    <Route
      {...rest}
      render={({ location }) =>
        auth.isLogin ? (
          children
        ) : (
          <Redirect
            to={{
              pathname: "/login",
              state: { from: location },
            }}
          />
        )
      }
    />
  );
}
