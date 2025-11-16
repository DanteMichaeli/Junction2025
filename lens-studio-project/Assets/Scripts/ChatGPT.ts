import { OpenAI } from "RemoteServiceGateway.lspkg/HostedExternal/OpenAI";

@component
export class ChatGPT extends BaseScriptComponent {
  @input internetModule: InternetModule;

  private ImageQuality = CompressionQuality.HighQuality;
  private ImageEncoding = EncodingType.Jpg;

  onAwake() {}

  makeImageRequest(imageTex: Texture, callback) {
    print("Making image request...");
    Base64.encodeTextureAsync(
      imageTex,
      (base64String) => {
        print("Image encode Success!");
        const textQuery =
          "Respond with ONLY one of these exact words: 'red-bull', 'vitamin-well-refresh', 'estrella-chips', or 'unknown'. Nothing else. Identify the main item in the image. Vitamin Well is a vitamin drink in a plastic bottle.";
        this.sendGPTChat(textQuery, base64String, callback);
      },
      () => {
        print("Image encoding failed!");
      },
      this.ImageQuality,
      this.ImageEncoding
    );
  }

  async sendGPTChat(request: string, image64: string, callback: (response: string) => void) {
    const gptStartTime = getTime();
    print("🕐 ChatGPT request started...");

    OpenAI.chatCompletions({
      model: "gpt-4o-mini",
      messages: [
        {
          role: "user",
          content: [
            { type: "text", text: request },
            {
              type: "image_url",
              image_url: {
                url: `data:image/jpeg;base64,` + image64,
              },
            },
          ],
        },
      ],
      max_tokens: 50,
    })
      .then((response) => {
        const gptEndTime = getTime();
        const gptDuration = (gptEndTime - gptStartTime) * 1000;
        print("⏱️  ChatGPT request completed in " + gptDuration.toFixed(2) + " ms");

        if (response.choices && response.choices.length > 0) {
          const content = response.choices[0].message.content;

          print("Content: " + content);

          if (content === "unknown") {
            callback("Unknown item or item not found");
          } else {
            this.sendToBackend(content);
            callback(`${content} added to cart!`);
          }
        }
      })
      .catch((error) => {
        const gptEndTime = getTime();
        const gptDuration = (gptEndTime - gptStartTime) * 1000;
        print("❌ ChatGPT request failed after " + gptDuration.toFixed(2) + " ms");
        print("Error in OpenAI request: " + error);
      });
  }

  async sendToBackend(itemName: string) {
    const backendStartTime = getTime();
    print("🕐 Backend request started...");
    print("Sending to Backend API: " + itemName);

    let request = new Request("https://af1d1fdc341f.ngrok-free.app/add-item-to-basket", {
      method: "POST",
      body: JSON.stringify({ itemId: itemName }),
      headers: {
        "Content-Type": "application/json",
      },
    });

    let response = await this.internetModule.fetch(request);
    const backendEndTime = getTime();
    const backendDuration = (backendEndTime - backendStartTime) * 1000;

    if (response.status != 201) {
      // API returns 201 Created
      print("❌ Backend request failed after " + backendDuration.toFixed(2) + " ms");
      print("Failure: response not successful. Status: " + response.status);
      return;
    }

    let contentTypeHeader = response.headers.get("Content-Type");
    if (!contentTypeHeader.includes("application/json")) {
      print("❌ Backend request failed after " + backendDuration.toFixed(2) + " ms");
      print("Failure: wrong content type in response");
      return;
    }

    let responseJson = await response.json();
    // The response structure is: { id, name, price }
    let itemId = responseJson.json["id"];
    let itemPrice = responseJson.json["price"];

    print("⏱️  Backend request completed in " + backendDuration.toFixed(2) + " ms");
    print("Item added successfully!");
    print("ID: " + itemId);
    print("Price: $" + itemPrice);

    return itemId;
  }
}
