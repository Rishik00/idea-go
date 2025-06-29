require('dotenv').config();
const { Client } = require('@notionhq/client');

const notion = new Client({ auth: process.env.NOTION_API_KEY });

(async () => {
  try {
    const response = await notion.pages.create({
      parent: {
        type: "page_id",
        page_id: "221f4dfa-12a3-801b-ae0a-cbd40522cbd4"
      },
      properties: {
        title: [
          {
            type: "text",
            text: {
              content: "💡 New Idea"
            }
          }
        ]
      },
      children: [
        {
          object: "block",
          type: "paragraph",
          paragraph: {
            rich_text: [
              {
                type: "text",
                text: {
                  content: "This is your idea description added from Node.js!"
                }
              }
            ]
          }
        }
      ]
    });

    console.log("✅ Page created:", response.url);
  } catch (error) {
    console.error("❌ Error creating page:", error.body || error);
  }
})();
