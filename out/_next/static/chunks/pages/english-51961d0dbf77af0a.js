(self.webpackChunk_N_E=self.webpackChunk_N_E||[]).push([[123],{6023:function(e,r,n){(window.__NEXT_P=window.__NEXT_P||[]).push(["/english",function(){return n(8050)}])},3193:function(e,r,n){"use strict";n.d(r,{f:function(){return o}});var t=n(5722),o={login:function(e,r){return t.JQ(r).then(n=>{console.log(r);var o=t.FH+"/api/user/member/login";return t.v_(o,{account:e,password:n})})}}},5722:function(e,r,n){"use strict";function t(e,r){return fetch(e,{method:"POST",body:JSON.stringify(r),headers:new Headers({"Content-Type":"application/json"})}).then(e=>e.json())}function o(e,r){var n=Date.now().toString();return l(n+a.sessionId).then(t=>fetch(e,{method:"POST",body:JSON.stringify(r),headers:new Headers({"Content-Type":"application/json",Authorization:[a.id,n,t].join(".")})}).then(e=>e.json()))}n.d(r,{EA:function(){return a},FH:function(){return s},JQ:function(){return l},T7:function(){return o},kS:function(){return c},ok:function(){return i},v_:function(){return t}});let s="https://pixie-383009.de.r.appspot.com",i=0;var a={id:"",sessionId:""};function l(e){let r=new TextEncoder().encode(e);return crypto.subtle.digest("SHA-256",r).then(e=>{let r=Array.from(new Uint8Array(e)),n=r.map(e=>e.toString(16).padStart(2,"0")).join("");return n})}function c(e){e.push("/")}},8050:function(e,r,n){"use strict";n.r(r),n.d(r,{default:function(){return B}});var t=n(5893),o=n(5722),s=n(7294),i=n(1163),a={ask:function(e){var r=o.FH+"/api/vocabulary/pool/ask";return o.T7(r,{userId:o.EA.id,vocabularyId:e})},back:function(){var e=o.FH+"/api/vocabulary/pool/back";return o.T7(e,{userId:o.EA.id})},forward:function(){var e=o.FH+"/api/vocabulary/pool/forward";return o.T7(e,{userId:o.EA.id})},jumpTo:function(e){var r=o.FH+"/api/vocabulary/pool/jump_to";return o.T7(r,{userId:o.EA.id,word:e})}},l={toggle:function(e){var r=o.FH+"/api/vocabulary/bookmark/toggle";return o.T7(r,{userId:o.EA.id,vocabularyId:e})},browse:function(e){var r=o.FH+"/api/vocabulary/bookmark/browse";return o.T7(r,{userId:o.EA.id,page:e})}};function c(e){var r=o.FH+"/api/sentence/interpret";return o.T7(r,{sentence:e})}var d=n(3619),u=n(1023),h=n(4096),x=n(9353),j=n(2864),p=n(7976),f=n(4385),Z=n(2085),g=n(2157),w=n(4789),v=n(2373),k=n(692),m=n(9879),S=n(7556),b=n(5913),C=n(911),y=n(9063),T=n(9789),I=n(9453),E=n(1772),_=n(9459),A=n(3414),D=n(3009),H=n(2424),P=n(6079),z={answer:{word:"",vocs:[],interpreteds:{}}};function F(){let[e,r]=(0,s.useState)(""),[n,H]=(0,s.useState)(""),[P,F]=(0,s.useState)(!1),[O,N]=(0,s.useState)([]),[J,L]=(0,s.useState)(!1),[U,Q]=(0,s.useState)(z.answer),[V,W]=(0,s.useState)(!1),[X,q]=(0,s.useState)(""),[B,G]=(0,s.useState)(!0),[K,M]=(0,s.useState)(""),[Y,$]=(0,s.useState)(!0),[ee,er]=(0,s.useState)(""),[en,et]=(0,s.useState)(!1),[eo,es]=(0,s.useState)(""),[ei,ea]=(0,s.useState)(""),el=(0,i.useRouter)();return(0,t.jsxs)(f.Z,{defaultValue:0,sx:{width:1e3,height:500,borderRadius:"lg",boxShadow:"sm",overflow:"auto"},children:[(0,t.jsxs)(Z.Z,{disableUnderline:!0,tabFlex:1,sx:{["& .".concat(g.Z.root)]:{fontSize:"sm",fontWeight:"lg",'&[aria-selected="true"]':{color:"primary.500",bgcolor:"background.surface"},["&.".concat(g.Z.focusVisible)]:{outlineOffset:"-4px"}}},children:[(0,t.jsx)(w.Z,{children:"Pool"}),(0,t.jsx)(w.Z,{children:"Unknown"}),(0,t.jsx)(w.Z,{children:"Explain"})]}),(0,t.jsx)(v.Z,{value:0,children:(0,t.jsxs)(k.Z,{container:!0,direction:"column",spacing:2,children:[(0,t.jsxs)(k.Z,{container:!0,direction:"row",justifyContent:"flex-start",children:[(0,t.jsx)(k.Z,{children:(0,t.jsxs)(m.Z,{size:"sm","aria-label":"soft button group",children:[(0,t.jsx)(S.Z,{onClick:e=>{e.preventDefault(),G(!0),a.back().then(e=>{if(e.errorCode!=o.ok){console.log(e),o.kS(el);return}G(!1),r(e.id),H(e.word),N(e.definitions),F(e.hasToggled)})},children:(0,t.jsx)(u.Z,{})}),(0,t.jsx)(S.Z,{onClick:e=>{e.preventDefault(),G(!0),a.forward().then(e=>{if(e.errorCode!=o.ok){console.log(e),o.kS(el);return}G(!1),r(e.id),H(e.word),N(e.definitions),F(e.hasToggled)})},children:(0,t.jsx)(d.Z,{})})]})}),(0,t.jsxs)(k.Z,{container:!0,direction:"row",alignItems:"center",children:[(0,t.jsx)(b.ZP,{size:"sm",value:X,placeholder:"enter word ...",onChange:e=>q(e.target.value),endDecorator:(0,t.jsx)(S.Z,{size:"sm",onClick:e=>{e.preventDefault(),a.jumpTo(X).then(e=>{if(e.errorCode!=o.ok){console.log(e),o.kS(el);return}if(q(""),!e.ok){G(!1);return}G(!0),r(e.id),H(e.word),N(e.definitions),F(e.hasToggled)})},children:"JumpTo"})}),(0,t.jsx)(C.ZP,{color:"danger",level:"body-xs",alignSelf:"flex-end",children:!B&&"The word is not found in pool"})]})]}),(0,t.jsxs)(k.Z,{container:!0,direction:"row",justifyContent:"center",spacing:2,children:[(0,t.jsx)(k.Z,{xs:5,sx:{height:"100%"},children:n&&(0,t.jsxs)(y.Z,{variant:"outlined",color:"primary",orientation:"horizontal",children:[(0,t.jsxs)(y.Z,{sx:{width:"100%",height:"100%",justifyContent:"center",alignItems:"flex-start"},children:[(0,t.jsx)(T.Z,{children:(0,t.jsx)(I.ZP,{size:"sm",color:P?"primary":"neutral",onClick:r=>{r.preventDefault(),l.toggle(e).then(e=>{if(e.errorCode!=o.ok){console.log(e),o.kS(el);return}F(!P)})},sx:{position:"absolute",zIndex:1,right:"1rem"},children:P?(0,t.jsx)(h.Z,{}):(0,t.jsx)(x.Z,{})})}),(0,t.jsx)(C.ZP,{level:"h1",children:n}),(0,t.jsx)(E.Z,{children:O.length>0&&O.map((e,r)=>(0,t.jsxs)(_.Z,{children:[e.partOfSpeech," |- ",e.text]},r))})]}),(0,t.jsx)(T.Z,{children:(0,t.jsx)(S.Z,{variant:"soft",loading:J,onClick:r=>{r.preventDefault(),L(!J),a.ask(e).then(e=>{if(e.errorCode!=o.ok){console.log(e),o.kS(el);return}z.answer.word=e.answer.word,z.answer.vocs=e.answer.vocs,z.answer.interpreteds={},Q(z.answer),L(!1)})},children:(0,t.jsx)(j.Z,{})})})]})}),(0,t.jsx)(k.Z,{xs:!0,children:U.vocs.length>0&&(0,t.jsx)(y.Z,{variant:"outlined",color:"warning",children:(0,t.jsx)(R,{answer:U,interpretCallback:(e,r)=>{W(!0),c(r).then(r=>{if(r.errorCode!=o.ok){o.kS(el);return}z.answer.interpreteds[e]=r.interpreted,Q(z.answer),W(!1)})}})})})]})]})}),(0,t.jsx)(v.Z,{value:1,children:(0,t.jsxs)(k.Z,{container:!0,direction:"column",justifyContent:"center",spacing:2,children:[(0,t.jsx)(k.Z,{container:!0,direction:"row",justifyContent:"center",alignItems:"center",children:(0,t.jsx)(b.ZP,{size:"sm",value:K,placeholder:"enter word ...",onChange:e=>M(e.target.value),endDecorator:(0,t.jsx)(S.Z,{size:"sm",onClick:e=>{var r;e.preventDefault(),$(!0),(r=o.FH+"/api/vocabulary/unknown/ask",o.T7(r,{userId:o.EA.id,word:K})).then(e=>{if(e.errorCode!=o.ok){o.kS(el);return}M(""),z.answer.word=e.answer.word,z.answer.vocs=e.answer.vocs,z.answer.interpreteds={},Q(z.answer),$(!1)})},children:"Ask"})})}),(0,t.jsx)(k.Z,{children:(0,t.jsx)(R,{answer:U,interpretCallback:(e,r)=>{W(!0),c(r).then(r=>{if(r.errorCode!=o.ok){o.kS(el);return}z.answer.interpreteds[e]=r.interpreted,Q(z.answer),W(!1)})}})})]})}),(0,t.jsx)(v.Z,{value:2,children:(0,t.jsxs)(k.Z,{container:!0,direction:"row",alignItems:"center",justifyContent:"space-between",children:[(0,t.jsx)(k.Z,{children:(0,t.jsx)(A.Z,{color:"primary",variant:"outlined",sx:{borderRadius:"6px"},children:(0,t.jsx)(D.Z,{value:eo,placeholder:"enter english text ...",required:!0,minRows:15,maxRows:15,sx:{width:400},onChange:e=>es(e.target.value)})})}),(0,t.jsx)(k.Z,{direction:"column",alignItems:"center",justifyContent:"center",children:(0,t.jsxs)(S.Z,{variant:"soft",loading:en,onClick:e=>{var r;e.preventDefault(),et(!0),(r=o.FH+"/api/sentence/explain",o.T7(r,{text:eo})).then(e=>{if(et(!1),e.errorCode!=o.ok){console.log(e),o.kS(el);return}ea(e.explained)})},children:["Explain",(0,t.jsx)(p.Z,{})]})}),(0,t.jsx)(k.Z,{children:(0,t.jsx)(A.Z,{color:"success",variant:"outlined",sx:{borderRadius:"6px"},children:(0,t.jsx)(D.Z,{placeholder:"explained ...",minRows:15,maxRows:15,value:ei,sx:{width:400}})})})]})})]})}function R(e){let{answer:r,interpretCallback:n}=e;return(0,t.jsx)(E.Z,{children:r.vocs.length>0&&r.vocs.map((e,o)=>(0,t.jsxs)(_.Z,{nested:!0,children:[(0,t.jsx)(H.Z,{sticky:!0,color:"warning",variant:"soft",children:e.word}),(0,t.jsxs)(E.Z,{marker:"disc",size:"sm",children:[(0,t.jsx)(_.Z,{color:"danger",children:e.partOfSpeech},"part_of_speech"),(0,t.jsx)(_.Z,{color:"primary",children:e.explain},"explain"),(0,t.jsx)(_.Z,{nested:!0,children:(0,t.jsx)(E.Z,{marker:"decimal",children:e.sentences.length>0&&e.sentences.map((e,s)=>(0,t.jsx)(_.Z,{children:(0,t.jsxs)(P.Z,{onClick:r=>{r.preventDefault(),n(o.toString()+s.toString(),e)},children:[e,r.interpreteds[o.toString()+s.toString()]]})},s))})},"sentences")]})]},o))})}var O=n(702),N=n(6457),J=n(1685);n(3193);var L=n(3540),U=n(2868),Q=n(9562),V=n(8717),W=n(5045),X=n(3950);function q(){let e=(0,i.useRouter)();return(0,t.jsxs)(U.L,{children:[(0,t.jsx)(Q.Z,{slots:{root:I.ZP},slotProps:{root:{variant:"plain",color:"neutral"}},sx:{borderRadius:40},children:(0,t.jsx)(V.Z,{children:o.EA.account})}),(0,t.jsx)(W.Z,{variant:"solid",invertedColors:!0,"aria-labelledby":"apps-menu-demo",sx:{"--List-padding":"0.5rem","--ListItemDecorator-size":"3rem",display:"grid",gridTemplateColumns:"repeat(3, 100px)",gridAutoRows:"100px",gap:1},children:(0,t.jsx)(X.Z,{orientation:"vertical",onClick:r=>{r.preventDefault(),o.kS(e)},children:(0,t.jsx)(L.Z,{})})})]})}function B(){let[e,r]=(0,s.useState)(1),[n,a]=(0,s.useState)([]),c=(0,i.useRouter)();return(0,t.jsxs)(k.Z,{container:!0,children:[(0,t.jsx)(k.Z,{container:!0,direction:"row",justifyContent:"center",children:(0,t.jsx)(q,{})}),(0,t.jsx)(k.Z,{children:(0,t.jsxs)(O.Z,{spacing:2,alignItems:"center",children:[(0,t.jsx)(N.Z,{children:"Learn"}),(0,t.jsx)(F,{}),(0,t.jsx)(N.Z,{children:"Bookmark"}),(0,t.jsxs)(m.Z,{children:[(0,t.jsx)(S.Z,{onClick:n=>{n.preventDefault(),l.browse(e-1).then(n=>{if(n.errorCode!=o.ok){console.log(n),o.kS(c);return}n.vocs.length>0&&(e>1&&r(e-1),a(n.vocs))})},children:(0,t.jsx)(u.Z,{})}),(0,t.jsx)(S.Z,{children:e}),(0,t.jsx)(S.Z,{onClick:n=>{n.preventDefault(),l.browse(e+1).then(n=>{if(n.errorCode!=o.ok){console.log(n),o.kS(c);return}n.vocs.length>0&&(r(e+1),a(n.vocs))})},children:(0,t.jsx)(d.Z,{})})]}),(0,t.jsxs)(J.Z,{children:[(0,t.jsx)("thead",{children:(0,t.jsxs)("tr",{children:[(0,t.jsx)("th",{children:"Id"}),(0,t.jsx)("th",{children:"Word"}),(0,t.jsx)("th",{children:"Part of speech"}),(0,t.jsx)("th",{children:"Text"})]})}),(0,t.jsx)("tbody",{children:n.length>0&&n.map((e,r)=>(0,t.jsxs)("tr",{children:[(0,t.jsx)("td",{children:e.id}),(0,t.jsx)("td",{children:e.word}),(0,t.jsx)("td",{children:e.partOfSpeech}),(0,t.jsx)("td",{children:e.text})]},r))})]})]})})]})}}},function(e){e.O(0,[774,895,178,888,179],function(){return e(e.s=6023)}),_N_E=e.O()}]);