require('dotenv').config();
const { Client } = require('@notionhq/client');

const notion = new Client({ auth: process.env.NOTION_API_KEY });

async function fetchChildren(blockId, depth = 0) {
  const indent = ' '.repeat(depth * 2);
  let cursor = undefined;

  do {
    const response = await notion.blocks.children.list({
      block_id: blockId,
      start_cursor: cursor,
    });

    for (const block of response.results) {
      console.log(`${indent}📄 Block ID: ${block.id}, Type: ${block.type}`);

      if (block.type === 'child_page') {
        console.log(`${indent}➡ Child Page Title: ${block.child_page.title}`);
      } else if (block.type === 'paragraph') {
        const textContent = block.paragraph.rich_text
          .map(t => t.plain_text)
          .join('');
        console.log(`${indent} ==> Paragraph: ${textContent}`);
      }

      if (block.has_children) {
        await fetchChildren(block.id, depth + 1);
      }
    }

    cursor = response.has_more ? response.next_cursor : null;
  } while (cursor);
}

(async () => {
  try {
    const blockId = process.env.PAGE_ID_ONE;
    console.log("📌 Root Block ID:", blockId);
    await fetchChildren(blockId);
  } catch (err) {
    console.error("❌ Error:", err);
  }
})();
