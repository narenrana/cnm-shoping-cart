import { createAsyncThunk, createSlice, PayloadAction } from "@reduxjs/toolkit";
import * as _ from "lodash";
import { httpClient } from "../../core";

const initialOtpState = {
  isLoading: false,
  discountCoupons: [],
};

const getCoupons = createAsyncThunk(
  "coupons/loadCoupons",
  async (request, thunkAPI) => {
    const options = {
      method: "GET",
      headers: {
        "Content-Type": "application/json",
      },
    };
    let response = (await httpClient("/coupons/v1/list", options)) || [];
    response = _.isEmpty(response) ? initialOtpState.defaultTabList : response;
    return response;
  }
);

const generateCoupons = createAsyncThunk(
  "coupons/generate",
  async (request, thunkAPI) => {
    console.log({ ...request.values });
    const options = {
      method: "POST",
      headers: {
        "Content-Type": "application/json",
      },
      body: JSON.stringify({ ...request }),
    };
    let response = (await httpClient("/coupons/v1/generate", options)) || [];
    response = _.isEmpty(response) ? initialOtpState.defaultTabList : response;
    return response;
  }
);
/**
 * Reducers
 */
const Reducer = createSlice({
  name: "coupons",
  initialState: initialOtpState,
  reducers: {},
  extraReducers: (builder) => {
    builder.addCase(getCoupons.fulfilled, (state, { payload = {} }) => {
      state.discountCoupons = (payload || {}).discountCoupons || [];
      state.isLoading = false;
    });

    builder.addCase(getCoupons.rejected, (state, action) => {
      state.isLoading = false;
    });
    builder.addCase(generateCoupons.fulfilled, (state, { payload }) => {
      state.discountCoupons = payload.discountCoupons;
      state.isLoading = false;
    });

    builder.addCase(generateCoupons.rejected, (state, action) => {
      state.isLoading = false;
    });
  },
});

/*Actions export*/
const {} = Reducer.actions;
export { generateCoupons, getCoupons };
/*Reducer export*/
export default Reducer.reducer;
