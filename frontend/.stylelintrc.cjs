module.exports = {
  plugins: [
    'stylelint-declaration-block-no-ignored-properties',
    'stylelint-order',
  ],
  extends: [
    './rules/stylelint/index.cjs',
  ],
  rules: {
    // 'order/properties-alphabetical-order': true,
    'order/properties-order': [
		[
			{
				emptyLineBefore: 'threshold',
        noEmptyLineBetween: true,
				properties: ['display', 'justify-content', 'align-items', 'flex-direction']
			},
			{
				emptyLineBefore: 'threshold',
        noEmptyLineBetween: true,
				properties: ['position', 'top', 'left', 'right', 'bottom', 'transform', 'z-index']
			},
			{
				emptyLineBefore: 'threshold',
        noEmptyLineBetween: true,
				properties: ['max-height', 'height', 'max-width', 'width', 'box-sizing']
			},
			{
				emptyLineBefore: 'threshold',
        noEmptyLineBetween: true,
				properties: ['padding', 'padding-top', 'padding-right', 'padding-bottom', 'padding-left', 'margin', 'margin-top', 'margin-right', 'margin-bottom', 'margin-left']
			},
			{
				emptyLineBefore: 'threshold',
        noEmptyLineBetween: true,
				properties: ['border', 'border-bottom', 'border-right', 'border-left', 'border-top', 'outline', 'border-radius']
			},
			{
				emptyLineBefore: 'threshold',
        noEmptyLineBetween: true,
				properties: ['font', 'font-size', 'font-family', 'font-weight', 'letter-spacing', 'line-height', 'text-align', 'text-decoration']
			},
			{
				emptyLineBefore: 'threshold',
        noEmptyLineBetween: true,
				properties: ['caret-color', 'color', 'background', 'background-color', 'opacity', 'box-shadow']
			},
			{
				emptyLineBefore: 'threshold',
        noEmptyLineBetween: true,
				properties: ['cursor', 'user-select', 'pointer-events']
			},
			{
				emptyLineBefore: 'threshold',
        noEmptyLineBetween: true,
				properties: ['animation', 'transition']
			},
		],
		{
			unspecified: 'bottom',
			emptyLineBeforeUnspecified: 'threshold',
			emptyLineMinimumPropertyThreshold: 4
		}
	]
  },
};
