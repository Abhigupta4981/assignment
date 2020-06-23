// For dev purposes
const dev = {
    s3: {
      REGION: "YOUR_BUCKET_REGION",
      BUCKET: "YOUR_BUCKET_NAME"
    },
    apiGateway: {
      REGION: "YOUR_API_GATEWAY_REGION",
      URL: "YOUR_API_GATEWAY_URL"
    },
    cognito: {
      REGION: "YOUR_COGNITO_REGION",
      USER_POOL_ID: "YOUR_COGNITO_USER_POOL_ID",
      APP_CLIENT_ID: "YOUR_COGNITO_APP_CLIENT_ID",
      IDENTITY_POOL_ID: "YOUR_COGNITO_IDENTITY_POOL_ID"
    }
};

// for production purposes
const prod = {
  s3: {
    REGION: "YOUR_BUCKET_REGION",
    BUCKET: "YOUR_BUCKET_NAME"
  },
  apiGateway: {
    REGION: "YOUR_API_GATEWAY_REGION",
    URL: "YOUR_API_GATEWAY_URL"
  },
  cognito: {
    REGION: "YOUR_COGNITO_REGION",
    USER_POOL_ID: "YOUR_COGNITO_USER_POOL_ID",
    APP_CLIENT_ID: "YOUR_COGNITO_APP_CLIENT_ID",
    IDENTITY_POOL_ID: "YOUR_COGNITO_IDENTITY_POOL_ID"
  }
};

const config = process.env.REACT_APP_STAGE === 'prod'? prod : dev;

export default config