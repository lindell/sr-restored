import { test, expect } from '@playwright/test';

test('basic navigation', async ({ page }) => {
	await page.goto('/');
	await page.getByRole('textbox', { name: 'Sök program' }).fill('dystopia');
	await page.getByRole('link', { name: 'P3 Dystopia' }).first().click();
	await expect(page.getByRole('heading', { level: 2, name: 'P3 Dystopia' })).toBeVisible();
});
