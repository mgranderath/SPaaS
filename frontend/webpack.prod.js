const HtmlWebpackPlugin = require('html-webpack-plugin');
const MiniCssExtractPlugin = require("mini-css-extract-plugin");

module.exports = {
  mode: 'production',
  resolve: {
    extensions: ['.js', '.vue']
  },
  module: {
    rules: [
      {
        test: /\.vue?$/,
        exclude: /(node_modules)/,
        use: 'vue-loader'
      },
      {
        test: /\.js?$/,
        exclude: /(node_modules)/,
        use: 'babel-loader'
      },
      {
        test: /\.css$/,
        use: [
          MiniCssExtractPlugin.loader,
          "css-loader", "postcss-loader",
        ],
      },
      {
        test: /\.svg$/,
        loader: 'vue-svg-loader', // `vue-svg` for webpack 1.x
      },
    ]
  },
  plugins: [
    new MiniCssExtractPlugin({
      filename: "index.css",
      chunkFilename: "index.css"
    }),
    new HtmlWebpackPlugin({
      template: './src/index.html'
    })
  ],
  devServer: {
    historyApiFallback: true
  },
  output: {
    publicPath: '/static',
  },
  externals: {
    // global app config object
    config: JSON.stringify({})
  }
};