# cheers-button

## Get started
1. `git clone git@github.com:kheiakiyama/chees-button.git`
2. Copy `terraform.tfvars.example` to `terraform.tfvars` and edit variables
3. `sh build_function.sh`
4. `terraform apply`
5. Set [AWS IoT Button](https://aws.amazon.com/iotbutton) call `send_webhook`.  
Each buttons, set placement attribute `room=meeting room A`,`room=meeting room B`,`room=meeting room C`...  

## Sample Payload
```
{
  "message": "Someone pushed the dash button just now. Want to join us for drinks at the meeting room A?",
  "meta": {
    "placementInfo": {
      "attributes": {
        "room": "meeting room A"
      }
    }
  }
}
```