# Automate Compliance
Provide a way to automate the compliance process for products, and turn them into consumable artifacts.

## Currently
* The NIST 800-53 product assessment is time consuming and manual in itself
* Translating the outcome of the assessment into OpenControl is currently a manual and error prone task
* Security compliance work is considered "boring/tedious", therefore it is difficult to get buy-in from product owners/SMEs
  * It doesn't help that the entire process is archane and manual work
* No metrics or insights into the overall process
* No accountability after the initial assessment

## Proposing
* Automate, Automate, Automate as much as possible
* Automate the process of translating the results from the 800-53 assessment into OpenControl
* Provide metrics to fill in the gaps and highlight areas of importance
* Consider integrating effort into other areas:
  * Creating STIG's/CIS or other standard baselines/benchmarks
  * Helping drive SCAP content development
  
The main goal is to drive automation throughout the process in order to make it easier for all parties involved. Achieving such a result would help shift compliance efforts to the left, and provide end consumers with quicker and more thoughtful outputs/artifacts to consume.

# Developer Setup
* Ensure you have the latest version of Go installed and configured (ie GOPATH)
* Follow **Step 1** of the instructions here: https://developers.google.com/sheets/api/quickstart/go
  * This will enable your Google Sheets API
  * Save your `credentials.json` to the root of this project
* `$ go run main.go`
* The first time you run the sample, it will prompt you to authorize access:
  * Browse to the provided URL in your web browser.
  * If you are not already logged into your Google account, you will be prompted to log in. If you are logged into multiple Google accounts, you will be asked to select one account to use for the authorization.
  * Click the Accept button.
  * Copy the code you're given, paste it into the command-line prompt, and press Enter.

