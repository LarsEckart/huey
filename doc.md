# Get Started
The fastest way to learn how to build apps which control the Hue system is to use the simple test web app built into every bridge. This lets you directly input commands and send them to the lights. You can look at the source HTML and JavaScript code for directions on how to do something different.
## Follow 3 Easy Steps
**Step 1**
First make sure your bridge is connected to your network and is functioning properly. Test that the smartphone app can control the lights on the same network.
**Step 2**
Discover the IP address of the bridge on your network. You can do this in a few ways. Note: When you are ready to make a production app, you should discover the bridge automatically using the Hue Bridge Discovery Guide.
1. Use an mDNS discovery app to find "Philips hue" on your network.
2. Use the broker server discover process by visiting https://discovery.meethue.com/
3. Log into your wireless router and look up "Philips hue" in the DHCP table.
4. Hue App method: Download the official Philips Hue app. Connect your phone to the network the Hue bridge is on. Start the Hue app. Push link connect to the bridge. Use the app to find the bridge and try controlling lights. If all works, go to the app settings → Hue Bridges → select your bridge to see its IP address.
**Step 3**
Once you have the address load the test app by visiting the following address in your web browser:
```
https://<bridge ip address>/debug/clip.html
```
You should see an interface that lets you populate the components of an HTTPS call — the basis of all web traffic and of the Hue RESTful interface.
- URL: the local address of a specific resource inside the Hue system (a light, group, etc.).
- Body: JSON describing what you want to change/add.
- Method: one of the HTTPS methods the Hue API uses:
- `GET` — fetch all information about the addressed resource
- `PUT` — modify an addressed resource
- `POST` — create a new resource inside the addressed resource
- `DELETE` — delete the addressed resource
- Response: JSON response from the bridge.
## So let’s get started…
First, get basic info about the bridge. Fill in the details leaving the body box empty and press `GET`.
URL:
```
/api/newdeveloper
```
Method:
```
GET
```
Next create a new authorized username. Fill in the info below and press `POST` (after pressing the bridge link button when prompted).
URL:
```
/api
```
Body:
```json
{"devicetype":"my_hue_app#iphone peter"}
```
Method:
```
POST
```
After pressing the bridge link button and submitting the `POST`, you should receive a success response with a username (for example: `1028d66426293e821ecfd9ef1a0731df`). Use that username for subsequent requests. Doing the first `GET` again will return much more information about lights and states in JSON format.
## Turning a light on and off
Get a list of all lights:
Address:
```
https://<bridge ip address>/api/1028d66426293e821ecfd9ef1a0731df/lights
```
Method:
```
GET
```
Get information about a specific light (example: light id 1):
Address:
```
https://<bridge ip address>/api/1028d66426293e821ecfd9ef1a0731df/lights/1
```
Method:
```
GET
```
Turn the light off (modify the `state` object):
Address:
```
https://<bridge ip address>/api/1028d66426293e821ecfd9ef1a0731df/lights/1/state
```
Body:
```json
{"on":false}
```
Method:
```
PUT
```
Turn the light on and change color/brightness/saturation:
Address:
```
https://<bridge ip address>/api/1028d66426293e821ecfd9ef1a0731df/lights/1/state
```
Body:
```json
{"on":true, "sat":254, "bri":254, "hue":10000}
```
Method:
```
PUT
```
You can vary the `hue` (0–65535), `sat` (0–254), and `bri` (0–254) values and re-send the `PUT` to see color changes.
Read more at the Core Concepts (developer account required): https://developers.meethue.com/develop/get-started-2/core-concepts/
